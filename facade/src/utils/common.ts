import dayjs from "dayjs"

export class JsonResponse {
    Data: any;
    Msg: string;
    Status: number;

    constructor(data: any, msg: string, status: number) {
        this.Data = data;
        this.Msg = msg;
        this.Status = status;
    }
}

export const timestamp2DateString = (timestamp: number) => {
    return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss")
}