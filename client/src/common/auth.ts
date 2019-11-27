import { useState, useCallback, createContext } from "react";
import { decode } from "jsonwebtoken";
import { UserInfo } from "interfaces/user";
import { AuthError, UnhandledError } from "interfaces/error";
import { signout } from "interfaces/endpoint";
import { useHttpTask } from "common/http";

type LogoutState = [undefined, null];
type LoginState = [undefined, UserInfo];
type InvalidState = [AuthError, undefined];
type AuthState = LogoutState | LoginState | InvalidState;

const LOGOUT_STATE: LogoutState = [undefined, null];

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
 * 引数のトークンからユーザ情報、エラー情報を返す
 */
const _login = (token: string): LoginState | InvalidState => {
  try {
    return [undefined, parse(token)];
  } catch (e) {
    if (e instanceof AuthError) {
      return [e, undefined];
    } else {
      throw new UnhandledError(e);
    }
  }
};

let initialState: LogoutState | LoginState = LOGOUT_STATE;
{
  let token = localStorage.getItem("jwt");
  if (token) {
    const result = _login(token);
    if (result[0]) {
      console.warn("logout caused by invalid jwt");
      localStorage.removeItem("jwt");
    } else {
      console.info("login initially");
      initialState = result;
    }
  }
}

type Handlers = [AuthState, (token: string) => void, () => void, () => void];

/**
 * ローカルストレージに保存されたJWTキーを解析し、ユーザ情報を取得する。
 * ```
 * [auth, login, logout, cancel] = useAuth();
 * login("new jwt");     // 再ログイン、ユーザ情報変更時
 * logout(null);         // ログアウト
 * ```
 * `auth` にユーザ情報が格納されている。
 * logout時サーバにリクエストを送る。中止させたいときは `cancel()` を呼び出す
 */
export const useAuth = (): Handlers => {
  const [auth, setAuth] = useState<AuthState>(initialState);
  const [logout, logoutCancel] = useHttpTask(signout, () => {
    console.info("task useAuth logout true");
    setAuth(LOGOUT_STATE);
    localStorage.removeItem("jwt");
  });

  const logoutWrapper = useCallback(() => logout(undefined), [logout]);
  const logoutCancelWrapper = useCallback(() => logoutCancel(undefined), [
    logout
  ]);

  const login = useCallback((t: string) => {
    console.info("callback useAuth login");
    const result = _login(t);
    setAuth(result);
    if (!result[0]) {
      console.info("save useAuth token");
      localStorage.setItem("jwt", t);
    }
  }, []);

  return [auth, login, logoutWrapper, logoutCancelWrapper];
};

/**
 * ```
 * const [auth, login, logout, cancel] = useContext(LoginContext)
 * ```
 * @see `useAuth`
 */
export default createContext<Handlers>([
  LOGOUT_STATE,
  () => {},
  () => {},
  () => {}
]);
