// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// #include "api.h"
#ifndef __BPF_API__
#define __BPF_API__

#include <stddef.h>
#include <linux/sched.h>
#include <linux/ptrace.h>
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>
#define _KERNEL(P)                                                                   \
	({                                                                     \
		typeof(P) val;                                                 \
		bpf_probe_read_kernel(&val, sizeof(val), &(P));                \
		val;                                                           \
	})

#define _(P)                                                                   \
	({                                                                     \
		typeof(P) val;                                                 \
		bpf_probe_read(&val, sizeof(val), &(P));                \
		val;                                                           \
	})

typedef enum
{
    true=1, false=0
} bool;

struct thread_struct {
    // x86_64
	long unsigned int fsbase;
	// arm64
	struct {
        unsigned long	tp_value;	/* TLS register */
        unsigned long	tp2_value;
    } uw;
}  __attribute__((preserve_access_index));

struct task_struct {
	__u32 pid;
    __u32 tgid;
    struct thread_struct thread;
}  __attribute__((preserve_access_index));
#endif

char __license[] SEC("license") = "Dual MIT/GPL";

struct key_t {
    __u32 user_stack_id;
    __u32 kernel_stack_id;
};

struct {
        __uint(type, BPF_MAP_TYPE_HASH);
        __type(key, struct key_t);
        __type(value, __u32);
        __uint(max_entries, 10000);
} counts SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_STACK_TRACE);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, 100 * sizeof(__u64));
    __uint(max_entries, 10000);
} stacks SEC(".maps");

SEC("perf_event")
int do_perf_event(struct pt_regs *ctx) {
    int monitor_pid;
    asm("%0 = MONITOR_PID ll" : "=r"(monitor_pid));

    // only match the same pid
    __u64 id = bpf_get_current_pid_tgid();
    __u32 tgid = id >> 32;

    if (tgid != monitor_pid) {
        return 0;
    }

    bpf_printk("flag 2");

    // create map key
    struct key_t key = {};

    // get stacks
    key.kernel_stack_id = bpf_get_stackid(ctx, &stacks, 0);
    key.user_stack_id = bpf_get_stackid(ctx, &stacks, BPF_F_USER_STACK);

    bpf_printk("kernel_stack_id: %d, user_stack_id: %d", key.kernel_stack_id, key.user_stack_id);

    __u32 *val;
    val = bpf_map_lookup_elem(&counts, &key);
    
    bpf_printk("&counts: %d", &counts);
    if (!val) {
        bpf_printk("flag 3");
        __u32 count = 0;
        bpf_map_update_elem(&counts, &key, &count, BPF_NOEXIST);
        val = bpf_map_lookup_elem(&counts, &key);
        if (!val) {
            bpf_printk("flag 4");
            return 0;
        }
        bpf_printk("flag 5");
    }
    bpf_printk("origin val: %d", *val);
    (*val) += 1;
    bpf_printk("update val: %d", *val);
    return 0;
}