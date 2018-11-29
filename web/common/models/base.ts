import { Monitorable } from "../interfaces/monitor";

const defaultValues: {id: string} = {id: "no value"};
const updateCallbacks: {[index: string]: (value: any) => void} = {};

/**
 * 変更監視ができるオブジェクト
 */
export default abstract class implements Monitorable {
    /**
     * 監視対象プロパティ
     */
    protected props: {[index: string]: any } = {};

    /**
     * 各プロパティが変更されたとき呼び出されるコールバック関数を格納した連想配列
     */
    protected updateCallbacks: {[index: string]: (value: any) => void} = {};

    /**
     * 監視開始前に呼び出すコールバック関数。
     * 配列の先頭から順に呼び出します。
     */
    protected beforeCallbacks: ((value: {[index: string]: any }) => void)[] = [];

    /**
     * 監視終了時に呼び出すコールバック関数。
     * 配列の後ろから順に呼び出します。
     */
    protected afterCallbacks: ((value: {[index: string]: any }) => void)[] = [];

    /**
     * 前回 reset時以降、値が更新されているかどうか
     */
    protected changed: boolean = true;

    /**
     * initialValueが存在しないときの値を設定します。
     * ここで登録した値は以降変更可能で監視対象になります。
     * @param props 
     */
    protected addDefaultValues(props: {[index: string]: {}} = {}) {
        Object.keys(props).forEach(key => this.props[key] = props[key]);
    }

    setupDefaultValues() {
        this.addDefaultValues(defaultValues);
    }

    setInitialValues(initialValues: {[index: string]: {}} = {}) {
        Object.keys(initialValues).filter(key => this.props[key] !== undefined)
            .forEach(key => this.props[key] = initialValues[key]);
    }

    /**
     * プロパティ更新時に呼び出すコールバック関数を設定します
     * @param key 監視対象のプロパティ名
     * @param callback 
     */
    protected addUpdateCallback(key: string, callback: (value: any) => void ) {
        this.updateCallbacks[key] = callback;
    }

    setupUpdateCallback() {
        Object.keys(updateCallbacks).forEach(key => {this.addUpdateCallback(key, updateCallbacks[key])});
    }

    setupBeforeCallback() {
        // do-nothing
    }

    begin() {
        this.beforeCallbacks.forEach(func => func(this.props));
    }

    setupAfterCallback() {
        // do-nothing
    }

    end() {
        // より基底のクラスのコールバックが最後に呼ばれるようにするため反転
        this.afterCallbacks.reverse().forEach(func => func(this.props));
    }

    /**
     * 対応するプロパティが定義されているとき、値を更新します
     * @param key プロパティ名
     * @param value プロパティ値
     */
    merge(key: string, value: any) {
        if (this.props[key] === undefined) {
            return;
        }

        if (this.props[key] != value) {
            this.props[key] = value;
            this.updateCallbacks[key](value);
            this.change();
        }
    }

    /**
     * propsに指定されたすべてのプロパティ値を更新します。
     * @param payload 
     */
    mergeAll(payload: {[index: string]:any} = {}) {
        Object.keys(payload).forEach((key => this.merge(key, payload[key])));
    }

    isChanged() {
        return this.changed;
    }

    change() {
        this.changed = true;
    }

    reset() {
        this.changed = false;
    }
}


