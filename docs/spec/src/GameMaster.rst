.. MIT License

    Copyright (c) 2017 yasshi2525

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

.. _gamemaster-spec:

ゲームマスタの仕様
==================

:term:`ゲームマスタ` は定期的に処理を実行して、ゲームの進行を行います。

.. _human-spec:

人の移動仕様
------------

:term:`人` が :term:`電車` を使って、 :term:`住宅` から :term:`会社` まで移動する仕様を定義します。

あらかじめ、すべて住宅と会社の組み合わせについて、経路を求めておきます。
これを :term:`移動テンプレート` と呼びます。

人が住宅で生成されるとき、該当するテンプレートを選択します。
これを :term:`人の移動経路` と呼びます。
テンプレート選択後は移動経路にしたがい行動します。

移動経路の変更
^^^^^^^^^^^^^^

移動中にマップが変化すると経路が変わる場合があります。
:term:`鉄道資産` の撤去により経路を辿れなくなったり、もっとよい経路が生まれる場合があります。

このため、マップに特定の変化があった際、再度経路探索を行います。
再計算時は、一旦すべての情報を破棄し、作成し直します。

.. note::

    実際に経路が変化するのは一部ですが、すべて再作成するのは、
    採用した経路探索アルゴリズムの仕様による制約です。


移動テンプレートを削除するとき、人の移動経路を削除します。
移動テンプレート再作成後、人の移動経路を個別に作成します。
(人の現在地点から会社までの経路は移動テンプレートにないため)

人は移動経路再計算中、自動行動をとります。

.. actdiag::

    actdiag {
        鉄道資産変更要求
        -> ゲーム進行停止
        -> 処理の中断
        -> 中間結果の破棄
        -> 鉄道資産変更
        -> 結果出力;

        鉄道資産変更
        -> ゲーム進行再開
        -> 人の移動経路削除
        -> 移動テンプレート削除
        -> 移動テンプレート作成
        -> 人の移動経路作成;

        // 人
        ゲーム進行再開
        -> 自由行動;

        人の移動経路作成
        -> 経路に沿って移動;

        lane プレイヤ {
            鉄道資産変更要求;
            結果出力;
        }

        lane ゲームマスタ {
            ゲーム進行停止;
            鉄道資産変更;
            人の移動経路削除;
            ゲーム進行再開;
        }

        lane 経路探索スレッド {
            処理の中断;
            移動テンプレート作成;
        }
        
        lane 人 {
            自由行動;
            経路に沿って移動;
        }
    }


計算中に再度マップに変化があった場合、計算を停止し、再計算します。

.. note::

    Ver 0.0では経路探索に関する情報をデータベースに永続化していましたが、
    大量のレコードが作成されたため、パフォーマンスの観点から Ver 1.0では永続化しません。
    経路情報はゲームマスタが管理します。アプリケーション停止すると経路情報は消去されるため、
    アプリケーション起動時に経路計算する機能を実装します。

経路探索の仕方について記述します。


.. _train-spec:

電車の走行仕様
--------------

:term:`電車` の走行アルゴリズムを説明します。

電車は :term:`路線` に所属し、路線の経路情報に従って、 :term:`線路` の上を走行します。

電車は以下の2つの状態があります。

#. 停車中
#. 走行中

停車を開始してから一定時間経過したら走り始めます。
走行していて駅についたら停車します。
路線で「通過」に設定されていれば停車せずに通過します。

特定の条件をみたすと、状態が変化します。

状態遷移図::
    
                  ---> 発車する --->
                  |               |
    待機する⇔ 1.停車中         2.走行中 ⇔ 走行する
                  |               |
                  <--- 停車する <---
    

待機する
^^^^^^^^

駅でドアを開いたまま待機します。
発車するまでのカウントを減らします。

電車の状態が停車中で、発車するまでのカウント(発車カウント)が1以上であれば待機します。

発車する
^^^^^^^^

ドアを閉じ、駅から発車します。移動はしません。
電車の状態を停車中から走行中に変更します。

電車の状態が停車中で、発車カウントが0であれば発車します。

走行する
^^^^^^^^

線路の上を移動します。
停車駅に到着する場合、オーバーランしないよう停車駅まで移動します。

電車が走行中で、停車駅に到達しなければ走行し続けます。

停車する
^^^^^^^^

駅に停車し、ドアを開きます。
電車の状態を走行中から停車中に変更します。
発車カウントをセットします。

電車が走行中で、停車駅に到着していれば停車します。

パラメタ
^^^^^^^^

発車カウントは、駅によって変わります。
駅の規模が大きいほど、停車する時間が長くなります。

線路を移動する距離は電車によって変わります。
電車の性能が良いほど、移動距離が長くなります。