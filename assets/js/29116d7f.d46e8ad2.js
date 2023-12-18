"use strict";(self.webpackChunkmgod=self.webpackChunkmgod||[]).push([[113],{2424:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>d,contentTitle:()=>c,default:()=>h,frontMatter:()=>t,metadata:()=>l,toc:()=>r});var i=s(5893),o=s(1151);const t={title:"Schema Options"},c=void 0,l={id:"schema_options",title:"Schema Options",description:"Schema Options is Mongo Schema level options (which modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.",source:"@site/../docs/schema_options.md",sourceDirName:".",slug:"/schema_options",permalink:"/mgod/docs/schema_options",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{title:"Schema Options"},sidebar:"docsSidebar",previous:{title:"Basic Usage",permalink:"/mgod/docs/basic_usage"},next:{title:"Field Options",permalink:"/mgod/docs/field_options"}},d={},r=[{value:"Collection",id:"collection",level:2},{value:"Usage",id:"usage",level:3},{value:"Timestamps",id:"timestamps",level:2},{value:"Usage",id:"usage-1",level:3},{value:"VersionKey",id:"versionkey",level:2},{value:"Usage",id:"usage-2",level:3},{value:"IsUnionType",id:"isuniontype",level:2},{value:"Usage",id:"usage-3",level:3},{value:"DiscriminatorKey",id:"discriminatorkey",level:2},{value:"Usage",id:"usage-4",level:3}];function a(e){const n={a:"a",admonition:"admonition",code:"code",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,o.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.p,{children:"Schema Options is Mongo Schema level options (which modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel."}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsx)(n.code,{children:"mgod"})," supports the following schema options -"]}),"\n",(0,i.jsx)(n.h2,{id:"collection",children:"Collection"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Accepts Type: ",(0,i.jsx)(n.code,{children:"string"})]}),"\n",(0,i.jsxs)(n.li,{children:["Is Optional: ",(0,i.jsx)(n.code,{children:"No"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["It is the name of the mongo collection in which the entity is stored. For example, ",(0,i.jsx)(n.code,{children:"users"})," collection of MongoDB for ",(0,i.jsx)(n.code,{children:"User"})," model in Golang."]}),"\n",(0,i.jsx)(n.h3,{id:"usage",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'schemaOpts := schemaopt.SchemaOptions{\n\tCollection: "users", // MongoDB collection name\n}\n'})}),"\n",(0,i.jsx)(n.h2,{id:"timestamps",children:"Timestamps"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Accepts Type: ",(0,i.jsx)(n.code,{children:"bool"})]}),"\n",(0,i.jsxs)(n.li,{children:["Default Value: ",(0,i.jsx)(n.code,{children:"false"})]}),"\n",(0,i.jsxs)(n.li,{children:["Is Optional: ",(0,i.jsx)(n.code,{children:"Yes"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["It is used to track ",(0,i.jsx)(n.code,{children:"createdAt"})," and ",(0,i.jsx)(n.code,{children:"updatedAt"})," meta fields for the entity. See ",(0,i.jsx)(n.a,{href:"/mgod/docs/meta_fields",children:"Meta Fields"})," for examples."]}),"\n",(0,i.jsx)(n.h3,{id:"usage-1",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'schemaOpts := schemaopt.SchemaOptions{\n\tCollection: "users",\n\tTimestamps: true,\n}\n'})}),"\n",(0,i.jsx)(n.h2,{id:"versionkey",children:"VersionKey"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Accepts Type: ",(0,i.jsx)(n.code,{children:"bool"})]}),"\n",(0,i.jsxs)(n.li,{children:["Default Value: ",(0,i.jsx)(n.code,{children:"true"})]}),"\n",(0,i.jsxs)(n.li,{children:["Is Optional: ",(0,i.jsx)(n.code,{children:"Yes"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["This reports whether to add a version key (",(0,i.jsx)(n.code,{children:"__v"}),") for the entity. See ",(0,i.jsx)(n.a,{href:"/mgod/docs/meta_fields",children:"Meta Fields"})," for examples."]}),"\n",(0,i.jsx)(n.h3,{id:"usage-2",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'schemaOpts := schemaopt.SchemaOptions{\n\tCollection: "users",\n\tVersionKey: true,\n}\n'})}),"\n",(0,i.jsx)(n.h2,{id:"isuniontype",children:"IsUnionType"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Accepts Type: ",(0,i.jsx)(n.code,{children:"bool"})]}),"\n",(0,i.jsxs)(n.li,{children:["Default Value: ",(0,i.jsx)(n.code,{children:"false"})]}),"\n",(0,i.jsxs)(n.li,{children:["Is Optional: ",(0,i.jsx)(n.code,{children:"Yes"})]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["It defines whether the entity is a union type. See ",(0,i.jsx)(n.a,{href:"/mgod/docs/union_types",children:"Union Types"})," for more details on unions."]}),"\n",(0,i.jsx)(n.h3,{id:"usage-3",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'schemaOpts := schemaopt.SchemaOptions{\n\tCollection: "resources",\n\tIsUnionType: true,\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["If ",(0,i.jsx)(n.code,{children:"IsUnionType"})," is set to true, then ",(0,i.jsx)(n.code,{children:"__t"})," will be used as the ",(0,i.jsx)(n.code,{children:"DiscriminatorKey"})," by default."]}),"\n",(0,i.jsx)(n.h2,{id:"discriminatorkey",children:"DiscriminatorKey"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Accepts Type: ",(0,i.jsx)(n.code,{children:"string"})]}),"\n",(0,i.jsxs)(n.li,{children:["Default Value: ",(0,i.jsx)(n.code,{children:"__t"})]}),"\n",(0,i.jsxs)(n.li,{children:["Is Optional: ",(0,i.jsx)(n.code,{children:"Yes"})]}),"\n"]}),"\n",(0,i.jsx)(n.p,{children:"It is the key used to identify the underlying type in case of a union type entity."}),"\n",(0,i.jsx)(n.h3,{id:"usage-4",children:"Usage"}),"\n",(0,i.jsx)(n.admonition,{type:"note",children:(0,i.jsxs)(n.p,{children:[(0,i.jsx)(n.code,{children:"IsUnionType"})," needs to be set to ",(0,i.jsx)(n.code,{children:"true"})," to use the ",(0,i.jsx)(n.code,{children:"DiscriminatorKey"})," field."]})}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'schemaOpts := schemaopt.SchemaOptions{\n\tCollection: "resources",\n\tIsUnionType: true,\n\tDiscriminatorKey: "type",\n}\n'})}),"\n",(0,i.jsx)(n.admonition,{type:"info",children:(0,i.jsxs)(n.p,{children:["The provided ",(0,i.jsx)(n.code,{children:"DiscriminatorKey"})," should be present in the Go struct as a compulsory field."]})}),"\n",(0,i.jsxs)(n.p,{children:["Default ",(0,i.jsx)(n.code,{children:"DiscriminatorKey"})," will be overwritten by the provided ",(0,i.jsx)(n.code,{children:"type"})," field."]})]})}function h(e={}){const{wrapper:n}={...(0,o.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(a,{...e})}):a(e)}},1151:(e,n,s)=>{s.d(n,{Z:()=>l,a:()=>c});var i=s(7294);const o={},t=i.createContext(o);function c(e){const n=i.useContext(t);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function l(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:c(e.components),i.createElement(t.Provider,{value:n},e.children)}}}]);