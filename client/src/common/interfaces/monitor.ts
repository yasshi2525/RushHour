import { Chunk, Point, ResolveError } from "./gamemap";

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
    updateDisplayInfo(): void

    /**
     * 監視を開始します。
     */
    begin(): void;

    get(key: string): any;
    position(): Point | undefined;

    /**
     * 引数のチャンク上に存在するか取得します。
     * @param chunk 
     */
    standOnChunk(chunk: Chunk): boolean;

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
     * 画面外の領域にあり、監視を終了すべきかどうかを返します。
     */
    shouldEnd(): boolean;

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

    resolve(error: ResolveError): ResolveError;
};

export interface MonitorContainer extends Monitorable {
    existsChild(id: string): boolean;
    getChild(id: string | undefined): Monitorable;
    getChildOnChunk(chunk: Chunk, oid: number): Monitorable | undefined
    mergeChild(payload: {id: string}): Monitorable;
    /**
     * このなかに存在しない child は outMap 属性が true になります
     * @param payload
     * @param opts 全child共通に設定するプロパティ
     */
    mergeChildren(payload: {id: string}[], opts: {[index: string]: any}): void;
    forEachChild(func: (c: Monitorable) => any): void;
    /**
     * shouldEndをみたすchildを削除します
     */
    endChildren(): void;
    removeChild(id: string): void;
}

export function startMonitor(model: Monitorable, props: {[index: string]: any}) {
    model.setupDefaultValues();
    model.setupUpdateCallback();
    model.setupBeforeCallback();
    model.setupAfterCallback();
    model.setInitialValues(props);
    model.begin();
}