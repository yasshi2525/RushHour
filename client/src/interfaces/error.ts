export enum ErrorType {
  /**
   * 認証に関するエラー(`401`)。再ログインを促す
   */
  AUTH,
  /**
   * サーバがメンテナス中によるエラー(`503`)。しらばく立ってからの操作を促す
   */
  OPERATION,
  /**
   * ユーザ操作が受け入れられなかったことによるエラー(`400`)。再操作を促す
   */
  REQUEST,
  /**
   * キャンセル命令によるエラー。再操作を促す
   */
  CANCEL,
  /**
   * サーバからの応答が不明だったことによるエラー。
   */
  SERVER,
  UNKNOWN,
  MULTI
}

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

  abstract get type(): ErrorType;

  abstract get summaries(): string[];
  get messages(): string[] {
    return this.reasons;
  }
}

/**
 * 認証に関するエラー(`401`)。再ログインを促す
 */
export class AuthError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.AUTH;
  }
  get summaries(): string[] {
    return ["認証エラーが発生しました", "再ログインしてください"];
  }
}
/**
 * サーバがメンテナス中によるエラー(`503`)。しらばく立ってからの操作を促す
 */
export class OperationError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.OPERATION;
  }
  get summaries(): string[] {
    return ["メンテナンス中です", "時間をおいてアクセスしてください"];
  }
}
/**
 * ユーザ操作が受け入れられなかったことによるエラー(`400`)。再操作を促す
 */
export class RequestError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.REQUEST;
  }
  get summaries(): string[] {
    return ["サーバへの要求が失敗しました", "再試行してください"];
  }
}

export class CancelError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.CANCEL;
  }
  get summaries(): string[] {
    return ["サーバへの要求がキャンセルされました", "再試行してください"];
  }
}

/**
 * サーバのレスポンスが不明だったことによるエラー。
 */
export class ServerError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.SERVER;
  }
  get summaries(): string[] {
    return ["サーバから不明な応答がありました", "再試行してください"];
  }
}

/**
 * アプリケーションが想定していないサーバ起因のエラー
 */
export class UnknownError extends GeneralError {
  get type(): ErrorType {
    return ErrorType.UNKNOWN;
  }
  get summaries(): string[] {
    return ["不明なエラーが発生しました", "画面を更新してください。"];
  }
}

export class MultiError extends GeneralError {
  private member: Errors[];

  get children(): Errors[] {
    return this.member;
  }

  constructor(errs: Errors[]) {
    super();
    this.member = errs;
  }

  get type(): ErrorType {
    return ErrorType.MULTI;
  }
  get summaries(): string[] {
    return ["複数の原因によるエラーが発生しました"];
  }
}

export const isMultiError = (obj: any): obj is MultiError =>
  obj instanceof Object && "children" in obj && obj.children instanceof Array;

export type Errors =
  | AuthError
  | OperationError
  | RequestError
  | ServerError
  | CancelError
  | UnknownError
  | MultiError;

export const isErrors = (obj: any): obj is Errors =>
  obj instanceof Object &&
  "type" in obj &&
  "summaries" in obj &&
  obj.summaries instanceof Array;

export const isOperationError = (obj: any): obj is OperationError =>
  isErrors(obj) && obj.type === ErrorType.OPERATION;
export const isAuthError = (obj: any): obj is AuthError =>
  isErrors(obj) && obj.type === ErrorType.AUTH;
export const isCancelError = (obj: any): obj is CancelError =>
  isErrors(obj) && obj.type === ErrorType.CANCEL;
