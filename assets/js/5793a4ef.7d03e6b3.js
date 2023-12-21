"use strict";(self.webpackChunkmgod=self.webpackChunkmgod||[]).push([[318],{8217:(n,e,o)=>{o.r(e),o.d(e,{assets:()=>c,contentTitle:()=>i,default:()=>m,frontMatter:()=>r,metadata:()=>a,toc:()=>d});var t=o(5893),s=o(1151);const r={title:"Transactions"},i=void 0,a={id:"transactions",title:"Transactions",description:"mgod provides a wrapper function WithTransaction that supports MongoDB transactions, allowing users to perform a series of read and write operations as a single atomic unit.",source:"@site/../docs/transactions.md",sourceDirName:".",slug:"/transactions",permalink:"/mgod/docs/transactions",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{title:"Transactions"}},c={},d=[{value:"Usage",id:"usage",level:2}];function l(n){const e={a:"a",admonition:"admonition",code:"code",h2:"h2",p:"p",pre:"pre",...(0,s.a)(),...n.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsxs)(e.p,{children:[(0,t.jsx)(e.code,{children:"mgod"})," provides a wrapper function ",(0,t.jsx)(e.code,{children:"WithTransaction"})," that supports MongoDB transactions, allowing users to perform a series of read and write operations as a single atomic unit."]}),"\n",(0,t.jsx)(e.h2,{id:"usage",children:"Usage"}),"\n",(0,t.jsxs)(e.p,{children:["Configure default connection with ",(0,t.jsx)(e.code,{children:"mgod"}),"."]}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}\nopts := options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=mgod_rs&authSource=admin")\n\nerr := mgod.ConfigureDefaultClient(cfg, opts)\n'})}),"\n",(0,t.jsx)(e.admonition,{type:"info",children:(0,t.jsxs)(e.p,{children:["To use Transactions, it is compulsory to run MongoDB daemon as a replica set.\nRefer Community Forum Discussion - ",(0,t.jsx)(e.a,{href:"https://www.mongodb.com/community/forums/t/why-replica-set-is-mandatory-for-transactions-in-mongodb/9533",children:"Why replica set is mandatory for transactions in MongoDB?"})]})}),"\n",(0,t.jsx)(e.p,{children:"Create models to be used inside a MongoDB transaction."}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'type User struct {\n\tName    string\n\tEmailID string `bson:"emailId"`\n}\n\ndbName := "mgoddb"\ncollection := "users"\nschemaOpts := schemaopt.SchemaOptions{\n\tTimestamps: true,\n}\n\nuserModel, _ := mgod.NewEntityMongoModelOptions(dbName, collection, &schemaOpts)\n'})}),"\n",(0,t.jsxs)(e.p,{children:["Use ",(0,t.jsx)(e.code,{children:"WithTransaction"})," function to perform multiple CRUD operations as an atomic unit."]}),"\n",(0,t.jsx)(e.pre,{children:(0,t.jsx)(e.code,{className:"language-go",children:'userDoc1 := User{Name: "Gopher1", EmailID: "gopher1@mgod.com"}\nuserDoc2 := User{Name: "Gopher2", EmailID: "gopher2@mgod.com"}\n\n_, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {\n\t_, err1 := s.userModel.InsertOne(sc, userDoc1)\n\t_, err2 := s.userModel.InsertOne(sc, userDoc2)\n\n\tif err1 != nil || err2 != nil {\n\t\treturn nil, errors.New("abort transaction")\n\t}\n\n\treturn nil, nil\n})\n'})}),"\n",(0,t.jsx)(e.admonition,{type:"warning",children:(0,t.jsxs)(e.p,{children:["Make sure to pass the session's context (",(0,t.jsx)(e.code,{children:"sc"})," here) only in EntityMongoModel's operation functions."]})})]})}function m(n={}){const{wrapper:e}={...(0,s.a)(),...n.components};return e?(0,t.jsx)(e,{...n,children:(0,t.jsx)(l,{...n})}):l(n)}},1151:(n,e,o)=>{o.d(e,{Z:()=>a,a:()=>i});var t=o(7294);const s={},r=t.createContext(s);function i(n){const e=t.useContext(r);return t.useMemo((function(){return"function"==typeof n?n(e):{...e,...n}}),[e,n])}function a(n){let e;return e=n.disableParentContext?"function"==typeof n.components?n.components(s):n.components||s:i(n.components),t.createElement(r.Provider,{value:e},n.children)}}}]);