import { Dispatch, createContext, useContext, useCallback } from "react";

export enum LoadingStatus {
  /**
   * ルートコンポーネントの作成完了
   */
  CREATED_ELEMENT,
  /**
   * 稼働中確認コンポーネントのロード完了
   */
  IMPORTED_OPERATION,
  /**
   * メンテナンス中か確認完了
   */
  CHECKED_OPERATION,
  /**
   * アプリケーションコンポーネントのロード完了
   */
  IMPORTED_APPLICATION,
  /**
   * メニューバーのインポート完了
   */
  IMPORTED_MENU,
  /**
   * ゲームボードのインポート完了
   */
  IMPORTED_BOARD,
  /**
   * 画像情報の読み込み完了
   */
  LOADED_RESOURCE,

  /**
   * プレイヤー情報の取得完了
   */
  FETCHED_PLAYERS,
  /**
   * マップ情報の取得完了
   */
  FETCHED_MAP,
  /**
   * キャンバスのロード完了
   */
  LOADED_CANVAS,
  /**
   * コントローラーの初期化完了
   */
  INITED_CONTROLLER,
  END
}

export namespace LoadingStatus {
  export function progress(st: LoadingStatus): number {
    return (st * 100) / LoadingStatus.END;
  }
  export function description(st: LoadingStatus): string {
    switch (st) {
      case LoadingStatus.CREATED_ELEMENT:
        return "ベースコンポーネントを読み込んでいます";
      case LoadingStatus.IMPORTED_OPERATION:
        return "ベースコンポーネントを構築しています";
      case LoadingStatus.CHECKED_OPERATION:
        return "アプリケーションコンポーネントを読み込んでいます";
      case LoadingStatus.IMPORTED_APPLICATION:
        return "メニューバーを読み込んでいます";
      case LoadingStatus.IMPORTED_MENU:
        return "ゲームコンポーネントを読み込んでいます";
      case LoadingStatus.IMPORTED_BOARD:
        return "画像情報を読み込んでいます";
      case LoadingStatus.LOADED_RESOURCE:
        return "プレイヤー情報を取得しています";
      case LoadingStatus.FETCHED_PLAYERS:
        return "マップ情報を取得しています";
      case LoadingStatus.FETCHED_MAP:
        return "コントローラーを構築しています";
      case LoadingStatus.LOADED_CANVAS:
        return "コントローラーを構築しています";
      case LoadingStatus.INITED_CONTROLLER:
        return "ロード処理を完了しています";
      default:
        console.warn(`invalid status : ${st}`);
        return `エラーが発生しました。画面をリロードしてください`;
    }
  }
}

type LoadingState = [LoadingStatus, Dispatch<LoadingStatus>];
const LoadingContext = createContext<LoadingState>([0, () => {}]);

export const useLoading = (): [LoadingStatus, Dispatch<LoadingStatus>] => {
  const [status, update] = useContext(LoadingContext);
  const updateWrapper = useCallback(
    (next: LoadingStatus) => {
      console.info(`callback useLoading ${status}=>${next}`);
      if (next > status) {
        update(next);
      } else {
        console.warn(`progress fallback: ${status}=>${next}`);
      }
    },
    [status, update]
  );
  return [status, updateWrapper];
};

export default LoadingContext;
