import React, { ReactNode, createContext } from "react";
import { decode } from "jsonwebtoken";
import { UserInfo } from "interfaces/user";
import { AuthError, Errors, isAuthError } from "interfaces/error";
import { signout } from "interfaces/endpoint";
import { OkResponse, ErrorResponse } from "./utils/http_common";
import { httpHttp } from "./utils/http_http";
import AsyncProvider from "./utils/provider_async";

const keySet = ["id", "name", "image", "admin", "hue"];
const fullKeySet = keySet.map(key => `${process.env.baseurl}/${key}`);

/**
 * 現在時刻と比較し、期限切れか判定する
 * @param exp 有効期限(UNIXタイムスタンプ(秒))
 */
const isExpired = (exp: any) => {
  if (exp === undefined || exp === null) {
    throw new AuthError(`jwt[exp] === ${exp}`);
  }
  const now = new Date().getTime() / 1000;
  return now > (exp as number);
};

/**
 * 取得したオブジェクトにすべてのキーが含まれているか判定する
 * @param obj `localStorage["jwt"]`のJSONをパースしたもの
 * @throws `AuthError` 必須キーが存在しないとき、値の型が一致しないとき
 */
const validate = (obj: { [index: string]: any }) => {
  const errs = fullKeySet
    .filter(key => !(key in obj))
    .map(key => `no jwt[${key}]`);
  if (errs.length > 0) {
    throw new AuthError(errs);
  }
  const result: { [index: string]: any } = {};
  keySet.forEach(key => {
    result[key] = obj[`${process.env.baseurl}/${key}`];
  });
  try {
    return result as UserInfo;
  } catch (e) {
    throw new AuthError(e);
  }
};

/**
 * JWTトークンからユーザ情報を取得する。
 * @throws `AuthError` 失敗時(必須キーがない、値の型が違う、期限切れ)
 */
const parse = (token: string) => {
  let obj = decode(token);
  if (obj == null || !(obj instanceof Object)) {
    throw new AuthError(`invalid jwt contents: ${obj}`);
  }
  if (isExpired(obj["exp"])) {
    throw new AuthError("jwt is expired");
  }
  return validate(obj);
};

/**
 * JWTの解析に成功すればローカルストレージに保存する。
 * 失敗した場合、ローカルストレージを削除する。
 * @throws `AuthError` 失敗時(必須キーがない、値の型が違う、期限切れ)
 */
const parseWithSync = (token: string) => {
  try {
    const result = parse(token);
    localStorage.setItem("jwt", token);
    return result;
  } catch (e) {
    if (isAuthError(e)) {
      localStorage.removeItem("jwt");
    }
    throw e;
  }
};

const authorize = (
  session: string | null
): OkResponse<string, UserInfo | null> | ErrorResponse<string> => {
  if (session == null) {
    return { args: "", payload: null };
  }
  try {
    return { args: session, payload: parseWithSync(session) };
  } catch (e) {
    if (e instanceof AuthError) {
      return { args: session, error: e };
    } else {
      throw e;
    }
  }
};

const initial = authorize(localStorage.getItem("jwt"));

type Endpoint = { args: string; payload: UserInfo | null };

const localEndpoint: Endpoint = {
  args: "",
  payload: null
};

const authTask = async (ep: Endpoint, signal: AbortSignal) => {
  if (ep.args) {
    return parseWithSync(ep.args);
  } else {
    await httpHttp(signout, signal);
    localStorage.removeItem("jwt");
    return null;
  }
};

const converter = (data: UserInfo | null) => data;

type AuthStatus = [
  UserInfo | null,
  (token: string) => void,
  (e: Errors) => void
];
const AuthContext = createContext<AuthStatus>([
  null,
  () => console.warn("not initialized AuthContext"),
  () => console.warn("not initialized AuthContext")
]);
AuthContext.displayName = "AuthContext";

interface ProviderProperties {
  children?: ReactNode;
  onError: ReactNode;
}

/**
 * ```
 * const [user, login, report] = useContext(AuthContext);
 * login("new_jwt_token");
 * login("") // sign out
 * ```
 */
export const AuthProvider = (props: ProviderProperties) => {
  return (
    <AsyncProvider
      ctx={AuthContext}
      initialFetch={initial}
      endpoint={localEndpoint}
      reloadTask={authTask}
      convert={converter}
      onError={props.onError}
    >
      {props.children}
    </AsyncProvider>
  );
};

export default AuthContext;
