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

.. _Glossary:

用語集
======

.. glossary::
    :sorted:

    閲覧者
        ログインをしていないユーザ。アカウントを持たない。
        マップを閲覧することができる。鉄道管理はできない。

    プレイヤ
        ゲームのプレイヤー。アカウントに紐づく。
        :term:`鉄道` を運営・管理できる。

    管理者
        サーバが稼働し続けるように管理する人。
        ゲームが進行するために、特権操作を行うことができる。

    ゲームマスタ
        ゲームの進行を担う主体。プレイヤではなく、ゲームスレッドのことを指す。
        タイマサービスによって、定期的に処理を実行する。

    マップ
        RushHourの都市。全 :term:`プレイヤ` の :term:`鉄道資産` と、その他資産を指す。
        その他資産とは、 :term:`住宅` 、 :term:`会社` 、 :term:`人` を指す。

    鉄道
        :term:`線路` 、:term:`駅` 、:term:`電車` 、:term:`路線` をひとまとめにしたもの。
        1 :term:`プレイヤ` 1鉄道所有することができる。

    鉄道資産
        :term:`プレイヤ` が所有する :term:`線路` 、:term:`駅` 、:term:`電車` 、:term:`路線` のこと。

    線路
        :term:`電車` が走るもの。上り線、下り線の2本の線から構成される。左側通行。

    駅
        :term:`電車` が停車するもの。
        :term:`人` が電車に乗り降りできる場所。
        
    プラットフォーム
        :term:`駅` の中にあり、 :term:`人` が :term:`電車` に乗り降りする場所。

    改札口
        :term:`駅` の中にあり、駅の外と :term:`プラットフォーム` に出入りできる場所。

    入場
        :term:`駅` の :term:`改札口` の外から :term:`プラットフォーム` に入ること。

    出場
        :term:`駅` の :term:`プラットフォーム` から :term:`改札口` の外に出ること。

    路線
        :term:`線路` と :term:`駅` から構成される経路情報。

    電車
        :term:`線路` の上を走るもの。
        1つの :term:`路線` に所属しており、路線の経路に従って線路の上を走る。

    人
        RushHourの住民。
        :term:`住宅` から :term:`会社` へ移動する。
        目的地までの移動経路は自分で決める。
        電車に乗ると :term:`乗客` になる。

    乗客
        電車に乗っている :term:`人` 。:term:`駅` で乗り降りする。

    住宅
        :term:`人` が住んでいる場所。
        ここから無限に人が生成され続ける。

    会社
        :term:`人` が勤めている場所。
        ここに到着すると人は消滅する。

    プレイヤトークン
        ログインした後、 :term:`プレイヤ` に対して発行する文字列。
        以降の操作はこの値を参照することで、どのプレイヤによる操作なのか判別する。
        
    ステップ
        ある地点からある地点まで移動すること。
        :term:`人` および :term:`電車` が移動する。
        :term:`電車` の場合、線路ノードから線路ノードまで移動することをステップと呼ぶ。

    移動テンプレート
        あらかじめ求められた、 :term:`住宅` から :term:`会社` までの移動経路。
        :term:`人用移動ステップ` が連続したもの

    人の移動経路
        人が選択した :term:`移動テンプレート` 。

    人用移動ステップ
        人がある点からある点まで移動すること。

    自由行動
        :term:`移動テンプレート` 再計算中に :term:`人` がとる行動。

    経路探索スレッド
        :term:`人` 用の経路探索を行うスレッド。
        :term:`ゲームマスタ` が作成・開始する。
        :term:`移動テンプレート` を出力する。