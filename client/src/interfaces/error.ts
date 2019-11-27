class GeneralError extends Error {
  private readonly reasons: string[];
  constructor(msg?: string | string[]) {
    if (msg === undefined) {
      super();
      this.reasons = [];
    } else if (msg instanceof Array) {
      super(msg.toString());
      this.reasons = msg;
    } else {
      super(msg);
      this.reasons = [msg];
    }
  }
  get messages(): Readonly<string[]> {
    return this.reasons;
  }
}

/**
 * 認証に関するエラー(`401`)。再ログインを促す
 */
export class AuthError extends GeneralError {}
/**
 * サーバがメンテナス中によるエラー(`503`)。しらばく立ってからの操作を促す
 */
export class OperationError extends GeneralError {}
/**
 * ユーザ操作が受け入れられなかったことによるエラー(`400`)。再操作を促す
 */
export class RequestError extends GeneralError {}
/**
 * アプリケーションが想定していないサーバ起因のエラー
 */
export class UnhandledError extends GeneralError {}

export type ServerErrors =
  | AuthError
  | OperationError
  | RequestError
  | UnhandledError;

export const isServerErrors = (obj: any): obj is ServerErrors => {
  return obj.messages instanceof Array;
};
