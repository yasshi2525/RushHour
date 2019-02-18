/**
 * 変更監視ができるクラスを表すインタフェース
 */
export interface Monitorable {
    /**
     * デフォルト値を設定します。
     * ここでセットされた値は監視対象プロパティになります。
     */
    setupDefaultValues(): void;

    /**
     * プロパティ更新時に呼び出すコールバック関数を設定します。
     */
    setupUpdateCallback(): void;

    /**
     * 初期値を代入します。
     * @param props 
     */
    setInitialValues(props: {[index: string]: {}}): void;

    /**
     * 変更監視開始前に呼び出されるコールバック関数を設定します。 
     */
    setupBeforeCallback(): void;

    /**
     * 変更監視終了前に呼び出されるコールバック関数を設定します。
     */
    setupAfterCallback(): void

    /**
     * 描画前に呼び出されるコールバック関数を設定します。
     */
    beforeRender(): void

    /**
     * 監視を開始します。
     */
    begin(): void;

    get(key: string): any;

    /**
     * keyに対応するプロパティが定義されているとき、値を更新します.
     * keyに対応するプロパティが定義されていないときは無視します。
     * @param key 
     * @param value 
     */
    merge(key:string, value: any): void;
    
    /**
     * 一括代入します。
     * @param payload プロパティが定義された連想配列
     */
    mergeAll(payload: {[index: string]: {}}): void; 

    /**
     * 監視を終了します。
     * @param callback 変更監視終了時に呼び出されるコールバック関数
     */
    end(): void;

    /**
     * 変更フラグをたてます
     */
    change(): void;

    /**
     * 前回の reset 以降に変更があったか取得します
     */
    isChanged(): boolean;

    /**
     * 変更フラグをおろします
     */
    reset(): void;
};

export interface MonitorContrainer extends Monitorable {
    existsChild(id: string): boolean;
    mergeChild(payload: {id: string}): void;
    /**
     * このなかに存在しない child は削除されます。
     * @param payload 
     */
    mergeChildren(payload: {id: string}[]): void;
    removeChild(id: string): void;
}