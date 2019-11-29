abstract class GeneralError extends Error {
  private reasons: string[];
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

  abstract get summary(): string;
  get messages(): string[] {
    return this.reasons;
  }
}

/**
 * 認証に関するエラー(`401`)。再ログインを促す
 */
export class AuthError extends GeneralError {
  get summary(): string {
    return "認証エラーが発生しました。再ログインしてください。";
  }
}
/**
 * サーバがメンテナス中によるエラー(`503`)。しらばく立ってからの操作を促す
 */
export class OperationError extends GeneralError {
  get summary(): string {
    return "メンテナンス中です。時間をおいてアクセスしてください。";
  }
}
/**
 * ユーザ操作が受け入れられなかったことによるエラー(`400`)。再操作を促す
 */
export class RequestError extends GeneralError {
  get summary(): string {
    return "エラーが発生しました。リトライしてください。";
  }
}
/**
 * アプリケーションが想定していないサーバ起因のエラー
 */
export class UnhandledError extends GeneralError {
  get summary(): string {
    return "不明なエラーが発生しました。";
  }
}

export type ServerErrors =
  | AuthError
  | OperationError
  | RequestError
  | UnhandledError;

export const isServerErrors = (obj: any): obj is ServerErrors =>
  obj instanceof Object && "summary" in obj;
