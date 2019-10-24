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

システム構成
============

システム全体図
--------------

RushHourはWebアプリケーションです。
外部システムとの関わりを定義します。

RushHour本体である、Webアプリケーション内の構成は :ref:`architecture-spec` で定義します。
図中の機能名の詳細は :ref:`function-spec` で定義します。

.. blockdiag::
    
    blockdiag {
        プレイヤ [shape = actor];
        管理者 [shape = actor];
        プレイヤ用Webサーバ [label = "Webサーバ"];
        管理者用Webサーバ [label = "Webサーバ"];
        認証サーバ [shape = "cloud"];
        データベース [shape = flowchart.database];
        タイマサービス [shape = roundedbox];

        認証機能 [shape = square];
        鉄道管理機能 [shape = square];
        マップ閲覧機能 [shape = square];
        ゲーム進行機能 [shape = square];
        ゲーム管理機能 [shape = square];


        group {
            label = "RushHour";
            fontsize = 16;
            認証機能, 鉄道管理機能, マップ閲覧機能, ゲーム進行機能, ゲーム管理機能
        }

        プレイヤ <-> プレイヤ用Webサーバ <-> 認証機能 <-> 認証サーバ;
                  プレイヤ用Webサーバ <-> 鉄道管理機能;
                  プレイヤ用Webサーバ <-> マップ閲覧機能;
        
        管理者 <-> 管理者用Webサーバ <-> ゲーム管理機能;
        ゲーム進行機能 <- タイマサービス;
        ゲーム管理機能 -> タイマサービス;

        // データベースとの関連矢印を表示すると、
        // 図が煩雑になるので、矢印表示をオフ
        認証機能         <-> データベース [style = none];
        鉄道管理機能     -> データベース [style = none];
        マップ閲覧機能   <- データベース [style = none];
        ゲーム進行機能   -> データベース [style = none];
        ゲーム管理機能   -> データベース [style = none];
    }

.. note ::

    各機能はデータベースとやりとりをしますが、簡略化のため関連線の表記を省略しています。

システム構成
------------

.. note ::

    以下は筆者yasshy2525がホスト ``rushhourgame.net`` で運用するときの構成です。
    将来的には、閲覧者が自身の環境にインストールできるようにしたいです。
    その際はRushHourが動作する環境を明記し、インストール手順ドキュメントを公開します。

.. list-table:: システム構成
    :header-rows: 1
    
    * - 要素
      - 

    * - OS
      - Cent OS 7

    * - Webサーバ
      - Nginx

    * - アプリケーションサーバ
      - Payara Server 4.1.1 以上

    * - DBMS
      - Maria DB 5.5 以上

    * - 使用言語
      - * Java
        * JavaScript

    * - バージョン管理
      - GitHub

    * - 依存性管理  
      - * Maven
        * npm

    * - ビルドツール
      - * Maven
        * gulp

    * - 継続的インテグレーションツール
      - Jenkins

    * - ドキュメント作成ツール
      - * Sphinx
        * javadoc

    * - Webアプリケーションフレームワーク
      - * Java Server Faces 2.2
        * PrimeFaces

    * - レンダリングエンジン
      - Pixi.js

    * - トランザクション管理
      - Java Transaction API 1.2

    * - O/R マッパ
      - EclipseLink 2.6.4

    * - 単体テストフレームワーク
      - * JUnit
        * Mockito (モック)
        * jacocco (カバレッジ測定)

    * - タスク管理・バグ管理
      - Redmine

    * - クライアントブラウザ
      - Google Chrome