(()=>{"use strict";var e,t,r,o,a,d={},n={};function f(e){var t=n[e];if(void 0!==t)return t.exports;var r=n[e]={id:e,loaded:!1,exports:{}};return d[e].call(r.exports,r,r.exports,f),r.loaded=!0,r.exports}f.m=d,f.c=n,e=[],f.O=(t,r,o,a)=>{if(!r){var d=1/0;for(u=0;u<e.length;u++){r=e[u][0],o=e[u][1],a=e[u][2];for(var n=!0,i=0;i<r.length;i++)(!1&a||d>=a)&&Object.keys(f.O).every((e=>f.O[e](r[i])))?r.splice(i--,1):(n=!1,a<d&&(d=a));if(n){e.splice(u--,1);var c=o();void 0!==c&&(t=c)}}return t}a=a||0;for(var u=e.length;u>0&&e[u-1][2]>a;u--)e[u]=e[u-1];e[u]=[r,o,a]},f.n=e=>{var t=e&&e.__esModule?()=>e.default:()=>e;return f.d(t,{a:t}),t},r=Object.getPrototypeOf?e=>Object.getPrototypeOf(e):e=>e.__proto__,f.t=function(e,o){if(1&o&&(e=this(e)),8&o)return e;if("object"==typeof e&&e){if(4&o&&e.__esModule)return e;if(16&o&&"function"==typeof e.then)return e}var a=Object.create(null);f.r(a);var d={};t=t||[null,r({}),r([]),r(r)];for(var n=2&o&&e;"object"==typeof n&&!~t.indexOf(n);n=r(n))Object.getOwnPropertyNames(n).forEach((t=>d[t]=()=>e[t]));return d.default=()=>e,f.d(a,d),a},f.d=(e,t)=>{for(var r in t)f.o(t,r)&&!f.o(e,r)&&Object.defineProperty(e,r,{enumerable:!0,get:t[r]})},f.f={},f.e=e=>Promise.all(Object.keys(f.f).reduce(((t,r)=>(f.f[r](e,t),t)),[])),f.u=e=>"assets/js/"+({53:"935f2afb",113:"29116d7f",199:"75b4a8fd",237:"1df93b7f",309:"726258db",346:"e745d665",368:"a94703ab",370:"fd90d755",388:"63b79d1b",518:"a7bd4aaa",616:"6dd8eb1d",623:"36c153ec",661:"5e95c892",705:"9895fb05",781:"3d529e1a",918:"17896441"}[e]||e)+"."+{53:"e1411b51",113:"d46e8ad2",199:"b2385d73",237:"942f6f59",309:"b0352525",346:"07fa5d9e",368:"6b0dd4cb",370:"67426e80",388:"2d8ac02e",518:"c8aea8b2",616:"7f627891",623:"9bf0bebc",661:"3a1a316a",705:"0b31c329",772:"dbefc956",781:"b5cf561a",918:"ef9fcff8"}[e]+".js",f.miniCssF=e=>{},f.g=function(){if("object"==typeof globalThis)return globalThis;try{return this||new Function("return this")()}catch(e){if("object"==typeof window)return window}}(),f.o=(e,t)=>Object.prototype.hasOwnProperty.call(e,t),o={},a="mgod:",f.l=(e,t,r,d)=>{if(o[e])o[e].push(t);else{var n,i;if(void 0!==r)for(var c=document.getElementsByTagName("script"),u=0;u<c.length;u++){var l=c[u];if(l.getAttribute("src")==e||l.getAttribute("data-webpack")==a+r){n=l;break}}n||(i=!0,(n=document.createElement("script")).charset="utf-8",n.timeout=120,f.nc&&n.setAttribute("nonce",f.nc),n.setAttribute("data-webpack",a+r),n.src=e),o[e]=[t];var b=(t,r)=>{n.onerror=n.onload=null,clearTimeout(s);var a=o[e];if(delete o[e],n.parentNode&&n.parentNode.removeChild(n),a&&a.forEach((e=>e(r))),t)return t(r)},s=setTimeout(b.bind(null,void 0,{type:"timeout",target:n}),12e4);n.onerror=b.bind(null,n.onerror),n.onload=b.bind(null,n.onload),i&&document.head.appendChild(n)}},f.r=e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},f.p="/mgod/",f.gca=function(e){return e={17896441:"918","935f2afb":"53","29116d7f":"113","75b4a8fd":"199","1df93b7f":"237","726258db":"309",e745d665:"346",a94703ab:"368",fd90d755:"370","63b79d1b":"388",a7bd4aaa:"518","6dd8eb1d":"616","36c153ec":"623","5e95c892":"661","9895fb05":"705","3d529e1a":"781"}[e]||e,f.p+f.u(e)},(()=>{var e={303:0,532:0};f.f.j=(t,r)=>{var o=f.o(e,t)?e[t]:void 0;if(0!==o)if(o)r.push(o[2]);else if(/^(303|532)$/.test(t))e[t]=0;else{var a=new Promise(((r,a)=>o=e[t]=[r,a]));r.push(o[2]=a);var d=f.p+f.u(t),n=new Error;f.l(d,(r=>{if(f.o(e,t)&&(0!==(o=e[t])&&(e[t]=void 0),o)){var a=r&&("load"===r.type?"missing":r.type),d=r&&r.target&&r.target.src;n.message="Loading chunk "+t+" failed.\n("+a+": "+d+")",n.name="ChunkLoadError",n.type=a,n.request=d,o[1](n)}}),"chunk-"+t,t)}},f.O.j=t=>0===e[t];var t=(t,r)=>{var o,a,d=r[0],n=r[1],i=r[2],c=0;if(d.some((t=>0!==e[t]))){for(o in n)f.o(n,o)&&(f.m[o]=n[o]);if(i)var u=i(f)}for(t&&t(r);c<d.length;c++)a=d[c],f.o(e,a)&&e[a]&&e[a][0](),e[a]=0;return f.O(u)},r=self.webpackChunkmgod=self.webpackChunkmgod||[];r.forEach(t.bind(null,0)),r.push=t.bind(null,r.push.bind(r))})()})();