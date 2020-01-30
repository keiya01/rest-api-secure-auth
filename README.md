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
  - **Cookie** 場合にもよるがこっちの方が良さそう?
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
  - 開発の段階で`secure`を`true`にしていると`localhost`で使用できない可能性があるため、開発時は`false`で良い(公開する時には`true`にすること)
- [x] `CORS`をちゃんと設定する([オリジン間リソース共有 (CORS)](https://developer.mozilla.org/ja/docs/Web/HTTP/CORS))
  - `Access-Control-Allow-Origin` ... 許可するOriginを指定する(デフォルトは同じOriginが指定される)
  - `Access-Control-Allow-Methods` ... 許可する`HTTP Method`を指定する(`GET, POST, OPTIONS, HEAD`など)
  - `Access-Control-Allow-Headers` ... 許可するヘッダーを指定する。プリフライトリクエストのレスポンスで使用される。(`Content-Type, Authorization`など)
  - `Access-Control-Allow-Credentials` ... 資格情報が必要なリクエストに対して、レスポンスを開示するかどうか(普通は何も指定しなくて良い) 
  - `Access-Control-Max-Age` ... プリフライトリクエストを何度も呼ぶのはオーバーヘッドになるので、このヘッダーに時間を指定することでキャッシュさせることができる
  - 上記の`CORS`をしっかり設定した上で`CSRF Token`をレスポンスする
- [x] `CSRF`対策をする
  - [gorilla/csrf](https://github.com/gorilla/csrf#javascript-applications)を使うと楽
  - CSRFの必要性([これで完璧！今さら振り返る CSRF 対策と同一オリジンポリシーの基礎](https://qiita.com/mpyw/items/0595f07736cfa5b1f50c), [gorilla/csrf で安全なWebフォームを作る](http://matope.hatenablog.com/entry/2019/06/05/144435))
  - cookieに CSRF Token を保存しておき、Client に Response する
  - Client では受け取った Token を Request に含めて送信する
  - `JWT`を使うことでステートレスなCSRF対策ができる(https://qiita.com/kaiinui/items/21ec7cc8a1130a1a103a)
- CSRF対策として`Preflight Request`を使う方法もあるが、`CSRF Token`を発行していれば、Same Origin であることの検証は可能なので必要ない
- [x] SQL Injection
- [x] Passwordなどの見られてはいけない重要な情報を暗号化してからDBに保存する
