# Rest API Secure Auth

# About
- golangでセキュアな認証システムを実装した
- 認証には goth を使っており、今回はTwitterログインのみを実装した
- SPAと組み合わせて使えるような REST API を意識して開発した
- 自動でログインされるようにcookieを保持している
- DB は Mock で作っており、メモリー上に情報を保持している
- 確認のためにここで確認できるようにカバーするべき項目を列挙している

# Secure API
- [x] 基本的に`gorilla`関連のpackageは様々な機能を小分けで提供してくれているため必要な物を組み合わせて安全に開発を行える
- [x] セッション情報は`Cookie`または`JWT`に保存する
  - **Cookie**
    - `Cookie`を使う場合は`httpOnly`や`secure`などのオプションから安全な設定を追加できるため、`XSS`からの攻撃を防ぐことできるが、`CSRF`は自前で実装して防ぐ必要がある
    - 安全に実装するなら`Cookie`のような気がする
  - **JWT**
    - `JWT`はブラウザの`LocalStrage`に保存することができ、`JWT`に情報を持たせることができるので`Server`をステートレスに保つことができる
    - `LocalStrage`は`Same-Origin`の場合のみでしか、I/O処理を行うことができないため`CSRF`の問題はないが、`XSS`によって情報を抜き取られる可能性がある
    - `XSS`を100%含まないと言い切れるサイトはない?(https://techracho.bpsinc.jp/hachi8833/2019_10_09/80851)
- [x] Cookie の扱いに気をつける
  - Cookie を信用しすぎない設計にする
  - ユーザー情報の編集などの個人情報の編集には必ず Password を求めるようにする
- [x] User ID を Cookie に保存する時は予測不可能なものに暗号化してからいれる(予測可能だとcookieを弄れば不正にログインできる)
  - `gorilla/sessions`を使うと楽
- [x] Cookieの`httpOnly`と`secure`を`true`にする(`httpOnly`はJSからアクセス不可能にするためで、`secure`は`https`でのみCookieを扱うことを指定する)
  - 開発の段階で`secure`を`true`にしていると localhost で使用できない可能性があるため、開発時は`false`で良い
- [x] CSRF対策をする(仕組みは[gorilla/csrf](https://github.com/gorilla/csrf#javascript-applications)または[gorilla/csrf で安全なWebフォームを作る](http://matope.hatenablog.com/entry/2019/06/05/144435)を見るとわかりやすい)
  - cookieに CSRF Token を保存しておき、Client に Response する
  - Client では受け取った Token を Request に含めて送信する
  - `JWT`を使うことでステートレスなCSRF対策ができる
- CSRF対策として`Preflight Request`もあるが、`CSRF Token`を発行していれば、Request Origin の検証は可能なので必要ないはず(間違っていたら教えてください、、)
- [x] SQL Injection
- [x] Passwordなどの見られてはいけない重要な情報を暗号化してからDBに保存する
