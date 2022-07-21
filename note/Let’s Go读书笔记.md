#### web应用的三大组件

- 处理器(handler)： responsible for executing your application logic and for writing HTTP response headers and bodies.
- 路由器(router)：stores a mapping between the URL patterns for your application and the corresponding handlers. Usually you have one servemux for your application containing all your routes

- web服务器(web server)

  

**Important: Go’s servemux treats the URL pattern "/" like a catch-all. So at the moment all HTTP requests will be handled by the home function regardless of their URL path. For instance, you can visit a different URL path like http://localhost:4000/foo, and you’ll receive exactly the same response. We’ll talk more about this in the next chapter.**



- Go’s servemux supports two different types of URL patterns: fixed paths and subtree paths. Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.

- fixed path patterns like these are only matched (and the corresponding handler called) when the request URL path exactly matches the fixed path.

- Subtree path patterns are matched (and the corresponding handler called) whenever the start of a request URL path matches the subtree path.

  

- In Go’s servemux, longer URL patterns always take precedence over shorter ones

- Request URL paths are automatically sanitized. If the request path contains any . or .. elements or repeated slashes, it will automatically redirect the user to an equivalent clean URL.

- If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a 301 Permanent Redirect to the subtree path with the slash added



- It’s only possible to call w.WriteHeader() once per response, and after the status code has been written it can’t be changed. If you try to call w.WriteHeader() a second time Go will log a warning message.

  w.WriteHeader()在一个handler中只能调用一次

- If you don’t call w.WriteHeader() explicitly, then the first call to w.Write() will automatically send a 200 OK status code to the user. So, if you want to send a non-200 status code, you must call w.WriteHeader() before any call to w.Write().

- change the header map after a call to w.WriteHeader() or w.Write() will have no effect on the response headers that the user receives. You need to make sure that your header map contains all the headers you want before you call these methods.

  ```go
  w.Header().Set("Allow", "POST") //该方法若放在WriteHeader()或Write()方法之后调用，则无效
  w.WriteHeader(405)
  w.Write([]byte("Method Not Allowed"))
  ```

- Header Canonicalization

  When you’re using the Add(), Get(), Set() and Del() methods on the header map, the header name will always be canonicalized using the textproto.CanonicalMIMEHeaderKey() function. This converts the first letter and any letter following a hyphen to upper case, and the rest of the letters to lowercase. This has the practical implication that when calling these methods the header name is case-insensitive. If you need to avoid this canonicalization behavior you can edit the underlying header map directly (it has the type map[string][]string). For example:

  ```go
  w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
  ```

- The sql.Open() function returns a sql.DB object. This isn’t a database connection — it’s a pool of many connections. This is an important difference to understand. Go manages these connections as needed, automatically opening and closing connections to the database via the driver.

- The connection pool is safe for concurrent access, so you can use it from web application handlers safely.

- The connection pool is intended to be long-lived. In a web application it’s normal to initialize the connection pool in your main() function and then pass the pool to your handlers. You shouldn’t call sql.Open() in a short-lived handler itself — it would be a waste of memory and network resources.

  

- Middleware：

  1. 中间件若作用于http.ServeMux，则该中间件对所有请求都生效
  2. 中间件若只作用于特定的用户自定义handler,则该中间件只对该特定的handler生效

  

- 任何为recover的panic都会导致应用程序崩溃推出
- go语言不支持跨协程recover程序触发的panic
- http请求所在的协程触发的panic会由go运行时自动recover，其他协程中的panic需要由用户手动recover