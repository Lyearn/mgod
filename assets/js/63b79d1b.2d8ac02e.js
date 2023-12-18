"use strict";(self.webpackChunkmgod=self.webpackChunkmgod||[]).push([[388],{7962:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>l,contentTitle:()=>d,default:()=>h,frontMatter:()=>s,metadata:()=>o,toc:()=>c});var i=r(5893),t=r(1151);const s={title:"Field Transformers"},d=void 0,o={id:"field_transformers",title:"Field Transformers",description:"Field transformers are an adapter between MongoDB field and Go struct field. They help in transforming field types in both directions i.e. from entity model to mongo doc and vice versa while building intermediate BSON document.",source:"@site/../docs/field_transformers.md",sourceDirName:".",slug:"/field_transformers",permalink:"/mgod/docs/field_transformers",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{title:"Field Transformers"},sidebar:"docsSidebar",previous:{title:"Field Options",permalink:"/mgod/docs/field_options"},next:{title:"Meta Fields",permalink:"/mgod/docs/meta_fields"}},l={},c=[{value:"ID",id:"id",level:2},{value:"Example",id:"example",level:3},{value:"Date",id:"date",level:2},{value:"Example",id:"example-1",level:3}];function a(e){const n={admonition:"admonition",code:"code",h2:"h2",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,t.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.p,{children:"Field transformers are an adapter between MongoDB field and Go struct field. They help in transforming field types in both directions i.e. from entity model to mongo doc and vice versa while building intermediate BSON document."}),"\n",(0,i.jsx)(n.admonition,{type:"note",children:(0,i.jsxs)(n.p,{children:["A field transformer is defined by the tag ",(0,i.jsx)(n.code,{children:"mgoType"}),"."]})}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsx)(n.code,{children:"mgod"})," supports the following field transformers -"]}),"\n",(0,i.jsx)(n.h2,{id:"id",children:"ID"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Tag Value: ",(0,i.jsx)(n.code,{children:"id"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["It is a transformer that converts a field of type ",(0,i.jsx)(n.code,{children:"string"})," in Go struct to ",(0,i.jsx)(n.code,{children:"primitive.ObjectID"})," for MongoDB document and vice versa."]}),"\n",(0,i.jsx)(n.h3,{id:"example",children:"Example"}),"\n",(0,i.jsx)(n.p,{children:"Type with id transformer."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n\tID   string `bson:"_id" mgoType:"id"`\n\tName string\n}\n\n// id = "65697705d4cbed00e8aba717"\nid := primitive.NewObjectID().Hex()\nuserDoc := User{\n\tID: id,\n\tName: "Gopher",\n}\n\nuser, _ := userModel.InsertOne(context.TODO(), userDoc)\n'})}),"\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Output:"})}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-js",children:'{\n\t"_id": ObjectId("65697705d4cbed00e8aba717"),\n\t"name": "Gopher"\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsx)(n.code,{children:"_id_"})," field will be of type ObjectId instead of String in MongoDB."]}),"\n",(0,i.jsx)(n.p,{children:"Invalid user doc -"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'userDoc := User{\n\tID: "randomId"\n\tName: "Gopher",\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["Inserting this doc will throw error as ",(0,i.jsx)(n.code,{children:"ID"})," field cannot be converted to primitive.ObjectID."]}),"\n",(0,i.jsx)(n.hr,{}),"\n",(0,i.jsx)(n.p,{children:"Type without id transformer."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n\tID   string\n\tName string\n}\n\nuserDoc := User{\n\tID: "randomId",\n\tName: "Gopher",\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["This is a valid doc now because there is no transformer applied on ",(0,i.jsx)(n.code,{children:"ID"})," field. Also, note that ",(0,i.jsx)(n.code,{children:"ID"})," field will be converted to ",(0,i.jsx)(n.code,{children:"id"})," instead of ",(0,i.jsx)(n.code,{children:"_id"})," because BSON tag is not present."]}),"\n",(0,i.jsx)(n.h2,{id:"date",children:"Date"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Tag Value: ",(0,i.jsx)(n.code,{children:"date"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["It is a transformer that converts a field of type ",(0,i.jsx)(n.code,{children:"string"})," in ISO 8601 format to ",(0,i.jsx)(n.code,{children:"primitive.DateTime"})," for MongoDB document and vice versa."]}),"\n",(0,i.jsx)(n.h3,{id:"example-1",children:"Example"}),"\n",(0,i.jsx)(n.p,{children:"Type with date transformer."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n\tName     string\n\tJoinedOn string `bson:"joinedOn" mgoType:"date"`\n}\n\n// joinedOn = "2023-12-01T11:32:19.290Z"\njoinedOn, _ := dateformatter.New(time.Now()).GetISOString()\nuserDoc := User{\n\tName: "Gopher",\n\tJoinedOn: joinedOn,\n}\n\nuser, _ := userModel.InsertOne(context.TODO(), userDoc)\n'})}),"\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Output:"})}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-js",children:'{\n\t"_id": ObjectId("65697705d4cbed00e8aba717"),\n\t"name": "Gopher",\n\t"joinedOn": ISODate("2023-12-01T11:32:19.290Z")\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsx)(n.code,{children:"joinedOn"})," field will be of type Date instead of String in MongoDB."]}),"\n",(0,i.jsx)(n.p,{children:"Invalid user doc -"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'userDoc := User{\n\tName: "Gopher",\n\tJoinedOn: "2023-12-01",\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["Inserting this doc will throw error as ",(0,i.jsx)(n.code,{children:"JoinedOn"})," field is not in expected ISO 8601 format."]}),"\n",(0,i.jsx)(n.hr,{}),"\n",(0,i.jsx)(n.p,{children:"Type without date transformer."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n\tName     string\n\tJoinedOn string `bson:"joinedOn"`\n}\n\nuserDoc := User{\n\tName: "Gopher",\n\tJoinedOn: "2023-12-01",\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["This is a valid doc now because there is no transformer applied on ",(0,i.jsx)(n.code,{children:"JoinedOn"})," field."]})]})}function h(e={}){const{wrapper:n}={...(0,t.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(a,{...e})}):a(e)}},1151:(e,n,r)=>{r.d(n,{Z:()=>o,a:()=>d});var i=r(7294);const t={},s=i.createContext(t);function d(e){const n=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:d(e.components),i.createElement(s.Provider,{value:n},e.children)}}}]);