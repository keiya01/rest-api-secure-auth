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
- [x] セッション情報はCookieに保存した方が安全(JWTはできるだけ避ける)
- [x] User ID を Cookie に保存する時は予測不可能なものに暗号化してからいれる(予測可能だとcookieを弄れば不正にログインできる)
  - `gorilla/sessions`を使うと楽
- [x] Cookieの`httpOnly`と`secure`を`true`にする(`httpOnly`はJSからアクセス不可能にするためで、`secure`は`https`でのみCookieを扱うことを指定する)
- [x] CSRF対策をする(仕組みは[gorilla/csrf](https://github.com/gorilla/csrf#javascript-applications)または[gorilla/csrf で安全なWebフォームを作る](http://matope.hatenablog.com/entry/2019/06/05/144435)を見るとわかりやすい)
- [x] SQL Injection(このRepoではDBを使っていない)
