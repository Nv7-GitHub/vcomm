export type Error = {
	message: string
}

export type OptionalError = {
	error?: Error
}

type SendData = {
	Fn: string,
	Value?: any
}

export type Sender = (value: string) => Promise<string>;

export type Param = {
	Val1: number;
	Val2: Array<string>;
}


export class Client extends EventTarget {
	sender: Sender;

	constructor(sender: Sender) {
		super();
		this.sender = sender;
	}

	async createMessage(fn: string, value?: any): Promise<any> {
		let data: SendData = {
			Fn: fn,
		}
		if (value) {
			data.Value = value;
		}
		let res = await this.sender(JSON.stringify(data));
		return JSON.parse(res);
	}

	async Receive(value: Param): Promise<OptionalError> {
		let res = await this.createMessage("Receive", value);
		return res as OptionalError;
	}
}

