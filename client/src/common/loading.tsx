import React, { Dispatch, createContext, useState, useCallback } from "react";
import CircularProgress from "@material-ui/core/CircularProgress";
import { ComponentProperty } from "interfaces/component";

export enum LoadingStatus {
  /**
   * ルートコンポーネントの作成完了
   */
  CREATED_ELEMENT,
  /**
   * 稼働中確認コンポーネントのロード完了
   */
  IMPORTED_OPERATION,
  CREATED_OPERATION,
  /**
   * メンテナンス中か確認完了
   */
  CHECKED_OPERATION,
  CREATED_MENU,
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
      case LoadingStatus.CREATED_OPERATION:
        return "ゲームの稼働ステータスを確認しています";
      case LoadingStatus.CHECKED_OPERATION:
        return "メニューバーを構築しています";
      case LoadingStatus.CREATED_MENU:
        return "画像データを読み込んでいます";
      case LoadingStatus.LOADED_RESOURCE:
        return "プレイヤー情報を取得しています";
      case LoadingStatus.FETCHED_PLAYERS:
        return "マップ情報を取得しています";
      case LoadingStatus.FETCHED_MAP:
        return "ハンドラを構築しています";
      case LoadingStatus.INITED_CONTROLLER:
        return "ロード処理を完了しています";
      default:
        console.warn(`invalid status : ${st}`);
        return `エラーが発生しました。画面をリロードしてください`;
    }
  }
}

interface LoadingState {
  status: LoadingStatus;
  update: Dispatch<LoadingStatus>;
}

const useLoading = (): LoadingState => {
  const [status, _update] = useState(0);
  const update = useCallback(
    (next: LoadingStatus) => {
      console.info(`callback useLoading ${status}=>${next}`);
      if (next > status) {
        _update(next);
      } else {
        console.warn(`progress fallback: ${status}=>${next}`);
      }
    },
    [status]
  );
  return { status, update };
};

export const LoadingCircle = () => (
  <CircularProgress aria-describedby="loading-description" aria-busy={true}>
    <div id="loading-description">読み込み中</div>
  </CircularProgress>
);

const LoadingContext = createContext<LoadingState>({
  status: 0,
  update: () => {}
});

export const LoadingProvider = (props: ComponentProperty) => {
  const context = useLoading();
  return (
    <LoadingContext.Provider value={context}>
      {props.children}
    </LoadingContext.Provider>
  );
};

export default LoadingContext;
