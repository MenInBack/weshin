微信公众号网页授权流程， weshin 层封装获取用户信息的过程。

```mermaid
sequenceDiagram
	participant user
	participant web
	participant weshin
	participant wechat
	
	user ->> web: 用户访问公众号页面
	
	web ->> web: 查询用户信息
	opt 已有用户信息
      web ->> user: web服务
	end
	
	web ->> weshin: 发起网页授权
	weshin ->> web: 授权跳转链接(jumpUrl)
	web -->> user: 引导页面跳转
	user ->> wechat: 跳转授权页面
	note over wechat: 用户授权
	wechat -->> user: 跳转回调页面
	user ->> web: 跳转 redirect_uri(code)
	web ->> weshin: 请求用户信息(code)
	weshin ->> wechat: 请求token (code|appID|secret)
	wechat ->> weshin: access_token|refresh_token
	opt optional
        note over weshin: 存储 token
    end
	weshin ->> wechat: 请求用户信息(access_token)
	wechat ->> weshin: 用户信息
    opt optional
        note over weshin: 存储用户信息
    end
	weshin ->> web: 用户信息
	web ->> user: web服务	
```

