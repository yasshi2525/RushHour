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

背景
====

ここでは、RushHourを開発しようと思った動機を説明します。

街があり、人々が生活している。
そこに、自分が何か操作すると、人々がそれを利用しだし、生活が変化する。
これをシミュレーションすることに興味があります。
もしかすると思いもよらない現象が見れるかもしれない点に面白さを感じているからです。

ただ、すべてをシミュレートするのはハードルが高すぎるので、
最も興味のある部分にフォーカスして、人々の生活→外部操作→生活の変化
のサイクルを実現したいと思いました。

そのエッセンスがRushHourです。

シミュレーション対象
""""""""""""""""""""

通勤時間帯の、人々が電車を使って一斉に会社に向かう様子をシミュレートする。
延々と広がる郊外から都心に向かう環境をイメージしています。
人々は様々な経路があるなかで、嗜好にあった行動をとります。

* 所要時間の短いルートで行く
* 運賃の安いルート [#fare]_ で行く
* 直通電車 [#direct]_ に乗って、乗り換え回数を少なくする
* 当駅始発電車 [#first]_ を待って座って行く
* 空いている各停 [#local]_ で行く

外部操作
"""""""""

以下の現象が起こったとします。

* 鉄道路線の新規開通
* 新規駅の開業
* 電車の増発
* 停車パターンの変更

人々の変化
"""""""""""

このような変化を期待しています。

* 新規路線に人が流れる
* 利便性の高いルートに変更


.. rubric:: 脚注

.. [#fare] 鉄道会社ごとの運賃の違いもあるが、実際は大差ない。
   それより、途中で他の鉄道会社の路線に乗り換える影響の方が大きい(初乗り運賃になってしまうため)。
   そのため、乗り換え回数を減らしたルートを選択する場合がある。

.. [#direct] ターミナルまで行く私鉄が、都心に直結する地下鉄と相互乗り入れしている。
   ターミナル付近で地下鉄に分岐し、そこで結構人が乗り換えている。

.. [#first] 途中駅でも車庫がある、折り返すなどで、始発電車が存在する。
   それを待って座って行く人はかなりおり、始発電車を待つ専用の列があったりする。

.. [#local] 急行運転を行っている路線を想定している。
   ターミナル駅に近い駅では、急行に乗り換えず、そのまま各停で乗り通す人が多い。
   短時間で乗り換えるのが面倒であったり、乗り換えても大した時間短縮にならなかったり、
   急行が混みすぎてそもそも乗れないなどの事情がある。
