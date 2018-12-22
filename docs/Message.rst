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

エラーメッセージ仕様
======================

ここではエラー時のメッセージ仕様を定義します。

エラーが発生したとき、エラー情報を下記の2つに対して出力します。

#. 画面
#. ログ

:term:`プレイヤ` は画面を通してエラー情報を取得します。
:term:`管理者` はログを通してエラー情報を取得します。

エラー情報を出力する目的は、プレイヤまたは管理者に対応を促すためです。

画面出力情報
------------

エラー発生時、以下の3種類の情報を画面に出力します。

#. エラーの概要
#. エラーの詳細
#. エラーの対処法

エラーの概要を示すことで、プレイヤが状況を理解しやすいようにします。

概要
^^^^

概要には以下の情報を出力します。

* プレイヤがしようとしていた動作
* その動作の結果

たとえば認証エラーが発生した場合は::

    ログインに失敗しました。

と出力します。プレイヤがしようとしていた動作は「ログインすること」で、
その動作の結果は「ログインに失敗した」です。

原因は詳細に出力します。「認証エラーが発生した」はログインに失敗した原因なので、
詳細に出力します。

詳細
^^^^

詳細には、エラーが発生した原因を出力します。
プレイヤに開示する意味のない情報は出力しません。
開示する意味のない情報とは、仮にプレイヤが知ったとしても
対処が打てない情報を指します。

たとえば認証エラーが発生した場合は::

    アカウントの認証に失敗しました。

と出力します。「認証サーバからの応答がなかった」
「認証サーバが不正な値を返した」は詳細に出力しません。

原因追求のための情報や、原因を取り除くための情報はログに出力します。

対処
^^^^

対処には、プレイヤが取るべき行動を出力します。
プレイヤがどのような行動を取ればエラーを回避できるか示します。

プレイヤの対処のしようがない場合は、
管理者への報告を促します。

.. note ::

    執筆時点ではどのようなエラーが発生するか不明なため、とりあえず
    「プレイヤが対処できないエラーが発生した場合、管理者へ報告する」
    という仕様にしました。
    そぐわない状況が判明した場合、この仕様を変更します。

再試行すれば成功できる可能性がある場合、再試行を促します。
たとえば認証サーバが一時的に応答不能であった場合、
時間を置いてアクセスすれば認証に成功する可能性があるので、再試行を促します。

たとえば認証エラーが発生した場合は::

    時間を置いて再度ログインして下さい。

であったり、::

    入力したアカウント情報が正しいか確認して下さい。

と出力します。

ログ出力情報
------------

ログには管理者が原因を特定したり、原因を取り除くのに必要な情報を出力します。

ログには以下の情報を出力します。

* エラーの発生箇所
* エラーの内容
* 入力として与えた値
* 実際の出力値

発生箇所
^^^^^^^^

どの処理でエラーが発生したか、一意に定まる情報を出力します。

エラーの内容
^^^^^^^^^^^^

期待に反して、何が発生したのか分かる情報を出力します。
以下の情報を含めたメッセージにして下さい。

* 何をしようとしていたか
* どうなることを期待していたか
* どのような状態になったか
* その結果、何ができなかったか

入力値
^^^^^^

プレイヤないしゲームマスタが与えた値を出力します。
ただし、セキュリティ上出力すべきでない値 [#secret]_ はマスキングします。

出力値
^^^^^^

計算の結果を出力します。
ただし、セキュリティ上出力すべきでない値 [#secret]_ はマスキングします。


.. rubric:: 脚注

.. [#secret] パスワード、認証キー、暗号化前の値など