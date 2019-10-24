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

.. _function-spec:

機能
====

RushHourがもつ機能について説明します。
大きく、プレイヤが操作する機能と、バックグラウンドで動く機能に大別されます。

ここでは以下のことを記載します。

* できること、できないこと
* 登場人物の説明 [#entity]_

以下のことは記載しません。

* どう実現するか
* 仕様 (○○の場合、できる/できない)

補足事項は本文が長くなるため、脚注に記載します。

.. rubric:: 脚注

.. [#entity] 何に対しての機能か分かるようにするため。

プレイヤができること
--------------------

* ログインする
* 自分の鉄道を管理する
* 経営による収支を管理する
* マップの様子を眺める
* 他プレイヤの経営状況を見る

認証機能
^^^^^^^^

プレイヤはアカウントを持ちます。
アカウントにログインすると、鉄道会社を経営できます。
ログインしていない場合は、マップの閲覧のみできます。

ログインの有無によるできることの違いを以下に示します。

.. csv-table:: 
    :header: ログイン,している,していない

    鉄道の管理,○,×         
    収支管理,○,×
    マップ閲覧,○,○
    経営状況閲覧,○,○ 

ログイン方法については :ref:`login-spec` で定義します。

鉄道管理機能
^^^^^^^^^^^^

鉄道は、線路、駅、電車の3つの物理的要素から構成されます。
電車は線路の上を走ります。乗客は駅で電車に乗り降りします。

論理的な要素として、路線が存在します。
電車は路線に従って運行します。

線路
""""

線路は複線とします [#double_track]_ 。

* 線路を敷設する
* 線路を撤去する

線路は途中で分岐させることができます。
分岐のコントロール方法は仕様の項目で記載します。

.. rubric:: 脚注

.. [#double_track] 運行しづらいため。正面衝突を回避する機能が必要になるため。

駅
""

* 駅を作る
* 駅を撤去する

電車
""""

* 電車を購入する
* 電車を線路に配置する
* 電車を線路から撤去する
* 電車を廃棄する

路線
""""

路線とは、始発駅から終着駅までの経路のことを指します。
途中駅の停車・通過を指定できます。

* 新規路線を作る
* 既存路線の経路を変更する
* 停車駅を変更する

路線機能を使って以下のことができます。

* 各駅停車と急行の2種別で運行する
* 長距離走る電車と、利用の多い区間だけ走る電車を運行する
* 途中駅でY字分岐させ、分岐・合流させる。

例::

    --- -1-> --- -2-> ---
    A駅      B駅      C駅
    --- <-4- --- <-3- ---

路線名「各停」::

    A駅停車 -> 1 -> B駅停車 -> 2 -> C駅停車 -> 3 -> B駅停車 -> 4 -> (最初に戻る)

路線名「急行」::

    A駅停車 -> 1 -> B駅通過 -> 2 -> C駅停車 -> 3 -> B駅通過 -> 4 -> (最初に戻る)

マップ閲覧機能
^^^^^^^^^^^^^^^^

.. todo :: 機能定義

経営状況確認機能
^^^^^^^^^^^^^^^^

.. todo :: 機能定義

管理者ができること
------------------

管理者はゲームの運用を維持するために以下のことができます。

* ゲーム時計の開始・停止
* 住宅、会社の建設・撤去

.. note ::

    以下の機能は検討中です。

    * ログインの凍結
    * プレイヤが作成した鉄道資産の削除

バックグラウンドで動く機能
--------------------------

ゲームを進行させるための機能です。
ゲームを進行させる主体をゲームマスタと呼びます。
ゲームマスタは人、電車を動かします。

ゲーム進行機能
^^^^^^^^^^^^^^

人が移動する
""""""""""""

人はRushHourの住民です。人は住宅から会社へ移動します。
移動の仕方は :ref:`human-spec` で定義します。

移動のモデルケースを以下に示します。

#. 住宅から生成される。
#. 駅まで徒歩で移動する。
#. 駅に入場する。
#. プラットフォームで電車を待つ。
#. 電車が到着したら、乗車する。
#. 目的地の最寄り駅についたら、下車する。
#. 駅から出場する。
#. 会社まで徒歩で移動する。
#. 会社に到着したら消滅する。

電車が走る
""""""""""

プレイヤが設置した電車を、路線で定義された経路に従って走行させます。
走り方は :ref:`train-spec` で定義します。

走行のモデルケースを以下に示します。

#. 駅から発車する
#. 線路の上を走行する
#. 駅についたら停車/通過する

電車はプレイヤが撤去しない限り、走行し続けます。


