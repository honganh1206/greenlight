# Web application workflow

If your API is a backend to a website rather than a standalone service, there are ways to *make it simpler and more intuitive for users* 

> [!WARNING]
> When creating a link in an email, do not rely on the Host header of `r.Host`, as that would lead to a [host header injection attack](https://portswigger.net/web-security/host-header)  

Method 1: Copy-and-paste the token into a form on your website, then perform the `PUT /v1/users/activated` request via JS - Simple and secure

```language
To activate your Greenlight account please visit h͟t͟t ͟p͟s͟:͟/͟/͟e͟x ͟a͟m͟p͟l͟e͟.͟c͟o͟m͟/͟u͟s͟e͟r͟s͟/͟a͟c͟t͟i͟v͟a͟t͟e͟
 and
enter the following code:
--------------------------
Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
--------------------------
```

Method 2: Ask the user to click on link (displayed as a button) which takes them to a page on your website

```language
To activate your Greenlight account please click the following link:
h͟t͟t͟p͟s͟:͟/͟/͟e͟x͟a͟m͟p͟l͟e͟.͟c͟o ͟m͟/͟u͟s͟e͟r͟s͟/͟a͟c͟t͟i͟v͟a͟t͟e͟?͟t͟o͟k͟e͟n͟=͟Y͟3͟Q͟M͟G͟X͟3͟P͟J͟3͟W ͟L͟R͟L͟2͟Y͟R͟T ͟Q͟G͟Q͟6͟K͟R͟H͟U͟
```

For the 2nd option, we need to avoid the token being [leaked in a referrer header](https://medium.com/@shahjerry33/password-reset-token-leak-via-referrer-2e622500c2c1) by setting `Referrer-Polcy: Origin`  
