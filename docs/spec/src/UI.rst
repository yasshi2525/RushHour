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

UI
==

:term:`プレイヤ` 、 :term:`管理者` のプレイヤインタフェース仕様を定義します。

.. _login-spec:

ログイン仕様
------------

プレイヤ
^^^^^^^^

アカウント登録の手間をプレイヤにかけさせたくないため、
外部サービスのアカウントを使ってログインできるようにします。

プレイヤは、以下のアカウントを使ってログインできます。

* Twitter

.. note ::

    OAuthが使えるサービスからログインできるよう設計します。

管理者
^^^^^^

管理者はパスワードを使ってログインします。

操作方針
--------

以下の2種類の方法で操作できるようにします。

* マップクリック
* コントローラ

マップをクリックすると、その地点でできることが表示されます。
再度クリックすると表示が消えます。

コントローラはマップと独立してあり、つねに表示されます。
ボタンが多いとビギナーが理解しづらくなるため、
可能な操作のみ表示します。
(駅が存在しないときは路線を組めないなど)

マップクリック
^^^^^^^^^^^^^^

マップをクリックすると、まずできる操作が列挙されます。
ユーザはその中からする操作を選びます。

.. list-table:: 各操作と表示される条件
    :header-rows: 1
    
    * - 操作名
      - 表示条件

    * - 線路開始 
      - * 線路ノードがない
        * 線路エッジがない

    * - 線路分割
      - 線路エッジがある

    * - 線路延伸
      - 線路ノードがある

    * - 線路削除
      - 線路ノードがある

    * - 駅設置
      - 線路ノードがある

    * - 駅編集
      - 駅がある

    * - 駅削除
      - 駅がある

    * - 路線作成
      - 駅がある

    * - 路線編集
      - 路線登録された駅がある

    * - 路線削除
      - 路線登録された駅がある

    * - 電車設置
      - 路線登録された駅がある

    * - 電車撤去
      - 電車がある