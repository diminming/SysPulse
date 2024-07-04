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