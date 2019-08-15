import { Monitorable } from "../interfaces/monitor";

const defaultValues: {id: number, oid: number, outMap: boolean} = {id: 0, oid: 0, outMap: false};

/**
 * 変更監視ができるオブジェクト
 */
export default class implements Monitorable {
    /**
     * 監視対象プロパティ
     */
    protected props: {[index: string]: any } = {};

    /**
     * 各プロパティが変更されたとき呼び出されるコールバック関数を格納した連想配列
     */
    protected updateCallbacks: {[index: string]: ((value: any) => void)[]} = {};

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
    protected changed: boolean = false;

    /**
     * initialValueが存在しないときの値を設定します。
     * ここで登録した値は以降変更可能で監視対象になります。
     * @param props 
     */
    addDefaultValues(props: {[index: string]: {}}) {
        Object.keys(props).forEach(key => {
            if (!(props[key] instanceof Object)) {
                this.props[key] = props[key]
            } else {
                this.props[key] = Object.assign({}, props[key]);
            }
        });
    }

    setupDefaultValues() {
        this.addDefaultValues(defaultValues);
    }

    setInitialValues(initialValues: {[index: string]: {}}) {
        Object.keys(initialValues).filter(key => this.props[key] !== undefined)
            .forEach(key => {
                if (!(initialValues[key] instanceof Object)) {
                    this.props[key] = initialValues[key]
                } else {
                    this.props[key] = Object.assign({}, initialValues[key]);
                }
            });
    }

    /**
     * プロパティ更新時に呼び出すコールバック関数を設定します
     * @param key 監視対象のプロパティ名
     * @param callback 
     */
    addUpdateCallback(key: string, callback: (value: any) => void ) {
        if (this.updateCallbacks[key] === undefined) {
            this.updateCallbacks[key] = [];
        }
        this.updateCallbacks[key].push(callback);
    }

    setupUpdateCallback() {
        //do-nothing
    }
    
    addBeforeCallback(handler: (value: {[index: string]: any }) => void) {
        this.beforeCallbacks.push(handler);
    }

    setupBeforeCallback() {
        // do-nothing
    }

    begin() {
        this.beforeCallbacks.forEach(func => func(this.props));
    }

    addAfterCallback(handler: (value: {[index: string]: any }) => void) {
        this.afterCallbacks.push(handler);
    }

    setupAfterCallback() {
        // do-nothing
    }

    beforeRender() {
        // do-nothing
    }

    shouldEnd() {
        return this.props.outMap;
    }

    end() {
        // より基底のクラスのコールバックが最後に呼ばれるようにするため反転
        this.afterCallbacks.reverse().forEach(func => func(this.props));
    }

    protected equals(key: string, value: any) {
        if (!(value instanceof Object)) {
            return this.props[key] == value;
        } else {
            return Object.keys(value).filter(k => this.props[key][k] != value[k]).length == 0;
        }
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

        if (!this.equals(key, value)) {
            if (!(value instanceof Object)) {
                this.props[key] = value;
            } else {
                this.props[key] = Object.assign({}, value)
            }
            if (this.updateCallbacks[key] !== undefined) {
                this.updateCallbacks[key].forEach(v => v(value));
            }
            this.change();
        }
    }

    /**
     * propsに指定されたすべてのプロパティ値を更新します。
     * @param payload 
     */
    mergeAll(payload: {[index: string]:any}) {
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

    get(key: string): any {
        return this.props[key];
    }
}


