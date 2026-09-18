import{Bt as e,D as t,E as n,Et as r,Gt as i,It as a,Jt as o,Kt as s,L as c,Lt as l,Mt as u,P as d,Qt as f,Rt as p,T as m,Ut as h,Zt as g,bt as _,d as v,gt as y,h as b,hn as x,ht as S,k as C,m as w,mn as T,nn as E,pn as D,rn as O,sn as k,vn as A,vt as j,xt as M,yt as N,z as P}from"./vue-i18n-BxOqVJSv.js";import{a as F,d as I,g as L,i as R,r as z,s as ee,t as te,u as ne}from"./create-CKNSf1FL.js";import{t as re}from"./format-length-CqnQVDgE.js";import{A as B,O as V,_ as ie,k as H}from"./_plugin-vue_export-helper-3sd4LcYN.js";import{c as ae,f as oe,g as se,h as U,p as ce,t as le,y as ue}from"./fade-in-scale-up.cssr-DpVYUTqJ.js";import{E as W,_ as de,p as fe,v as pe}from"./index-0G3VD5Ph.js";function me(e={},t){let n=D({ctrl:!1,command:!1,win:!1,shift:!1,tab:!1}),{keydown:r,keyup:i}=e,a=e=>{switch(e.key){case`Control`:n.ctrl=!0;break;case`Meta`:n.command=!0,n.win=!0;break;case`Shift`:n.shift=!0;break;case`Tab`:n.tab=!0}r!==void 0&&Object.keys(r).forEach(t=>{if(t!==e.key)return;let n=r[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},o=e=>{switch(e.key){case`Control`:n.ctrl=!1;break;case`Meta`:n.command=!1,n.win=!1;break;case`Shift`:n.shift=!1;break;case`Tab`:n.tab=!1}i!==void 0&&Object.keys(i).forEach(t=>{if(t!==e.key)return;let n=i[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},s=()=>{(t===void 0||t.value)&&(B(`keydown`,document,a),B(`keyup`,document,o)),t!==void 0&&k(t,e=>{e?(B(`keydown`,document,a),B(`keyup`,document,o)):(H(`keydown`,document,a),H(`keyup`,document,o))})};return ce()?(g(s),f(()=>{(t===void 0||t.value)&&(H(`keydown`,document,a),H(`keyup`,document,o))})):s(),T(n)}var G=h({name:`ChevronRight`,render(){return(()=>{let e=m(`6ab04425f4fcb756`);return e[0]||=l(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[l(`path`,{d:`M5.64645 3.14645C5.45118 3.34171 5.45118 3.65829 5.64645 3.85355L9.79289 8L5.64645 12.1464C5.45118 12.3417 5.45118 12.6583 5.64645 12.8536C5.84171 13.0488 6.15829 13.0488 6.35355 12.8536L10.8536 8.35355C11.0488 8.15829 11.0488 7.84171 10.8536 7.64645L6.35355 3.14645C6.15829 2.95118 5.84171 2.95118 5.64645 3.14645Z`,fill:`currentColor`})],-1)})()}});function K(e){return t=>{e.value=t?t.$el:null}}var he={...R,...w.props},ge=h({name:`Tooltip`,props:he,slots:Object,__popover__:!0,setup(e){let{mergedClsPrefixRef:t}=d(e),n=w(`Tooltip`,`-tooltip`,void 0,de,e,t),r=x(null);return{syncPosition(){r.value.syncPosition()},setShow(e){r.value.setShow(e)},popoverRef:r,mergedTheme:n,popoverThemeOverrides:a(()=>n.value.self)}},render(){let{mergedTheme:e,internalExtraClass:t}=this;return i(z,{...this.$props,theme:e.peers.Popover,themeOverrides:e.peerOverrides.Popover,builtinThemeOverrides:this.popoverThemeOverrides,internalExtraClass:t.concat(`tooltip`),ref:`popoverRef`},this.$slots)}}),_e=y(`icon`,`
 height: 1em;
 width: 1em;
 line-height: 1em;
 text-align: center;
 display: inline-block;
 position: relative;
 fill: currentColor;
`,[N(`color-transition`,{transition:`color .3s var(--n-bezier)`}),N(`depth`,{color:`var(--n-color)`},[S(`svg`,{opacity:`var(--n-opacity)`,transition:`opacity .3s var(--n-bezier)`})]),S(`svg`,{height:`1em`,width:`1em`})]),ve={...w.props,depth:[String,Number],size:[Number,String],color:String,component:[Object,Function]},ye=h({_n_icon__:!0,name:`Icon`,inheritAttrs:!1,props:ve,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=d(e),r=w(`Icon`,`-icon`,_e,fe,e,t),i=a(()=>{let{depth:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value;if(t!==void 0){let{color:e,[`opacity${t}Depth`]:r}=i;return{"--n-bezier":n,"--n-color":e,"--n-opacity":r}}return{"--n-bezier":n,"--n-color":``,"--n-opacity":``}}),o=n?b(`icon`,a(()=>`${e.depth||`d`}`),i,e):void 0;return{mergedClsPrefix:t,mergedStyle:a(()=>{let{size:t,color:n}=e;return{fontSize:re(t),color:n}}),cssVars:n?void 0:i,themeClass:o?.themeClass,onRender:o?.onRender}},render(){let{$parent:e,depth:t,mergedClsPrefix:n,component:r,onRender:a,themeClass:s}=this;return e?.$options?._n_icon__&&P(`icon`,"don't wrap `n-icon` inside `n-icon`"),a?.(),i(`i`,o(this.$attrs,{role:`img`,class:[`${n}-icon`,s,{[`${n}-icon--depth`]:t,[`${n}-icon--color-transition`]:t!==void 0}],style:[this.cssVars,this.mergedStyle]}),r?i(r):this.$slots.default?.())}}),q=c(`n-dropdown-menu`),J=c(`n-dropdown`),Y=c(`n-dropdown-option`),X=h({name:`DropdownDivider`,props:{clsPrefix:{type:String,required:!0}},render(){return E(),e(`div`,{class:n(`${this.clsPrefix}-dropdown-divider`)},null,2)}});function Z(e,t){return e.type===`submenu`||e.type===void 0&&e[t]!==void 0}function be(e){return e.type===`group`}function Q(e){return e.type===`divider`}function xe(e){return e.type===`render`}function Se(e,t,n){if(!t)return e;let r=x(e.value),i=null;return k(e,e=>{i!==null&&window.clearTimeout(i),e===!0?n&&!n.value?r.value=!0:i=window.setTimeout(()=>{r.value=!0},t):r.value=!1}),r}var $=h({name:`DropdownOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null},placement:{type:String,default:`right-start`},props:Object,scrollable:Boolean},setup(e){let t=s(J),{hoverKeyRef:n,keyboardKeyRef:r,lastToggledSubmenuKeyRef:i,pendingKeyPathRef:o,activeKeyPathRef:c,animatedRef:l,mergedShowRef:u,renderLabelRef:d,renderIconRef:f,labelFieldRef:p,childrenFieldRef:m,renderOptionRef:h,nodePropsRef:g,menuPropsRef:_}=t,y=s(Y,null),b=s(q),S=s(U),C=a(()=>e.tmNode.rawNode),w=a(()=>{let{value:t}=m;return Z(e.tmNode.rawNode,t)}),T=a(()=>{let{disabled:t}=e.tmNode;return t}),E=Se(a(()=>{if(!w.value)return!1;let{key:t,disabled:a}=e.tmNode;if(a)return!1;let{value:s}=n,{value:c}=r,{value:l}=i,{value:u}=o;return s===null?c===null?l!==null&&u.includes(t):u.includes(t)&&u[u.length-1]!==t:u.includes(t)}),300,a(()=>r.value===null&&!l.value)),D=a(()=>!!y?.enteringSubmenuRef.value),k=x(!1);O(Y,{enteringSubmenuRef:k});function A(){k.value=!0}function j(){k.value=!1}function M(){let{parentKey:t,tmNode:a}=e;a.disabled||u.value&&(i.value=t,r.value=null,n.value=a.key)}function N(){let{tmNode:t}=e;t.disabled||u.value&&n.value!==t.key&&M()}function P(t){if(e.tmNode.disabled||!u.value)return;let{relatedTarget:r}=t;r&&!L({target:r},`dropdownOption`)&&!L({target:r},`scrollbarRail`)&&(n.value=null)}function F(){let{value:n}=w,{tmNode:r}=e;u.value&&!n&&!r.disabled&&(t.doSelect(r.key,r.rawNode),t.doUpdateShow(!1))}return{labelField:p,renderLabel:d,renderIcon:f,siblingHasIcon:b.showIconRef,siblingHasSubmenu:b.hasSubmenuRef,menuProps:_,popoverBody:S,animated:l,mergedShowSubmenu:a(()=>E.value&&!D.value),rawNode:C,hasSubmenu:w,pending:v(()=>{let{value:t}=o,{key:n}=e.tmNode;return t.includes(n)}),childActive:v(()=>{let{value:t}=c,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r!==-1&&r<t.length-1}),active:v(()=>{let{value:t}=c,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r!==-1&&r===t.length-1}),mergedDisabled:T,renderOption:h,nodeProps:g,handleClick:F,handleMouseMove:N,handleMouseEnter:M,handleMouseLeave:P,handleSubmenuBeforeEnter:A,handleSubmenuAfterEnter:j}},render(){let{animated:a,rawNode:s,mergedShowSubmenu:c,clsPrefix:l,siblingHasIcon:d,siblingHasSubmenu:f,renderLabel:m,renderIcon:h,renderOption:g,nodeProps:_,props:v,scrollable:y}=this,b=null;if(c){let e=this.menuProps?.(s,s.children);b=(t=>(E(),p(Ee,o({key:1},e,{clsPrefix:l,scrollable:this.scrollable,tmNodes:this.tmNode.children,parentKey:this.tmNode.key}),null,16,[`clsPrefix`,`scrollable`,`tmNodes`,`parentKey`])))(b)}let x={class:[`${l}-dropdown-option-body`,this.pending&&`${l}-dropdown-option-body--pending`,this.active&&`${l}-dropdown-option-body--active`,this.childActive&&`${l}-dropdown-option-body--child-active`,this.mergedDisabled&&`${l}-dropdown-option-body--disabled`],onMousemove:this.handleMouseMove,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onClick:this.handleClick},S=_?.(s),w=(E(),e(`div`,o({class:[`${l}-dropdown-option`,S?.class],"data-dropdown-option":!0},S),[C(()=>i(`div`,o(x,v),[(E(),e(`div`,{class:n([`${l}-dropdown-option-body__prefix`,d&&`${l}-dropdown-option-body__prefix--show-icon`])},[C(()=>[h?h(s):W(s.icon)])],2)),(E(),e(`div`,{"data-dropdown-option":!0,class:n(`${l}-dropdown-option-body__label`)},[m?(E(),e(u,{key:0},[C(()=>m(s))],64)):(E(),e(u,{key:1},[C(()=>W(s[this.labelField]??s.title))],64))],2)),(E(),e(`div`,{"data-dropdown-option":!0,class:n([`${l}-dropdown-option-body__suffix`,f&&`${l}-dropdown-option-body__suffix--has-submenu`])},[this.hasSubmenu?(E(),p(ye,{key:0},{_:1,default:t(()=>(E(),p(G)))})):C(()=>null)],2))])),this.hasSubmenu?(E(),p(I,{key:0},{default:()=>[(E(),p(ne,null,{default:()=>(E(),e(`div`,{class:n(`${l}-dropdown-offset-container`)},[(E(),p(ee,{show:this.mergedShowSubmenu,placement:this.placement,to:y&&this.popoverBody||void 0,teleportDisabled:!y},{default:()=>(E(),e(`div`,{class:n(`${l}-dropdown-menu-wrapper`)},[a?(E(),p(r,{key:0,onBeforeEnter:this.handleSubmenuBeforeEnter,onAfterEnter:this.handleSubmenuAfterEnter,name:`fade-in-scale-up-transition`,appear:!0},{default:()=>b},1032,[`onBeforeEnter`,`onAfterEnter`])):(E(),e(u,{key:1},[C(()=>b)],64))],2))},1032,[`show`,`placement`,`to`,`teleportDisabled`]))],2))},1024))]},1024)):C(()=>null)],16));return g?g({node:w,option:s}):w}}),Ce=h({name:`DropdownGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{showIconRef:e,hasSubmenuRef:t}=s(q),{renderLabelRef:n,labelFieldRef:r,nodePropsRef:i,renderOptionRef:a}=s(J);return{labelField:r,showIcon:e,hasSubmenu:t,renderLabel:n,nodeProps:i,renderOption:a}},render(){let{clsPrefix:t,hasSubmenu:r,showIcon:i,nodeProps:a,renderLabel:s,renderOption:c}=this,{rawNode:d}=this.tmNode,f=(E(),e(`div`,o({class:`${t}-dropdown-option`},a?.(d)),[l(`div`,{class:n(`${t}-dropdown-option-body ${t}-dropdown-option-body--group`)},[l(`div`,{"data-dropdown-option":!0,class:n([`${t}-dropdown-option-body__prefix`,i&&`${t}-dropdown-option-body__prefix--show-icon`])},[C(()=>W(d.icon))],2),l(`div`,{class:n(`${t}-dropdown-option-body__label`),"data-dropdown-option":!0},[s?(E(),e(u,{key:0},[C(()=>s(d))],64)):(E(),e(u,{key:1},[C(()=>W(d.title??d[this.labelField]))],64))],2),l(`div`,{class:n([`${t}-dropdown-option-body__suffix`,r&&`${t}-dropdown-option-body__suffix--has-submenu`]),"data-dropdown-option":!0},null,2)],2)],16));return c?c({node:f,option:d}):f}}),we=h({name:`NDropdownGroup`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null}},render(){let{tmNode:t,parentKey:n,clsPrefix:r}=this,{children:a}=t;return E(),e(u,null,[(E(),p(Ce,{clsPrefix:r,tmNode:t,key:t.key},null,8,[`clsPrefix`,`tmNode`])),C(()=>a?.map(e=>{let{rawNode:t}=e;return t.show===!1?null:Q(t)?i(X,{clsPrefix:r,key:e.key}):e.isGroup?(P(`dropdown`,"`group` node is not allowed to be put in `group` node."),null):(E(),p($,{clsPrefix:r,tmNode:e,parentKey:n,key:e.key},null,8,[`clsPrefix`,`tmNode`,`parentKey`]))}))],64)}}),Te=h({name:`DropdownRenderOption`,props:{tmNode:{type:Object,required:!0}},render(){let{rawNode:{render:e,props:t}}=this.tmNode;return i(`div`,t,[e?.()])}}),Ee=h({name:`DropdownMenu`,props:{scrollable:Boolean,showArrow:Boolean,arrowStyle:[String,Object],clsPrefix:{type:String,required:!0},tmNodes:{type:Array,default:()=>[]},parentKey:{type:[String,Number],default:null}},setup(e){let{renderIconRef:t,childrenFieldRef:n}=s(J);O(q,{showIconRef:a(()=>{let n=t.value;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>n?n(e):e.icon);let{rawNode:t}=e;return n?n(t):t.icon})}),hasSubmenuRef:a(()=>{let{value:t}=n;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>Z(e,t));let{rawNode:n}=e;return Z(n,t)})})});let r=x(null);return O(se,null),O(ue,null),O(U,r),{bodyRef:r}},render(){let{parentKey:t,clsPrefix:r,scrollable:i}=this,a=this.tmNodes.map(e=>{let{rawNode:n}=e;return n.show===!1?null:xe(n)?(E(),p(Te,{tmNode:e,key:e.key},null,8,[`tmNode`])):Q(n)?(E(),p(X,{clsPrefix:r,key:e.key},null,8,[`clsPrefix`])):be(n)?(E(),p(we,{clsPrefix:r,tmNode:e,parentKey:t,key:e.key},null,8,[`clsPrefix`,`tmNode`,`parentKey`])):(E(),p($,{clsPrefix:r,tmNode:e,parentKey:t,key:e.key,props:n.props,scrollable:i},null,8,[`clsPrefix`,`tmNode`,`parentKey`,`props`,`scrollable`]))});return E(),e(`div`,{class:n([`${r}-dropdown-menu`,i&&`${r}-dropdown-menu--scrollable`]),ref:`bodyRef`},[i?(E(),p(ie,{key:0,contentClass:`${r}-dropdown-menu__content`},{default:()=>a},1032,[`contentClass`])):(E(),e(u,{key:1},[C(()=>a)],64)),this.showArrow?(E(),e(u,{key:2},[C(()=>F({clsPrefix:r,arrowStyle:this.arrowStyle,arrowClass:void 0,arrowWrapperClass:void 0,arrowWrapperStyle:void 0}))],64)):C(()=>null)],2)}}),De=y(`dropdown-menu`,`
 transform-origin: var(--v-transform-origin);
 background-color: var(--n-color);
 border-radius: var(--n-border-radius);
 box-shadow: var(--n-box-shadow);
 position: relative;
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
`,[le(),y(`dropdown-option`,`
 position: relative;
 `,[S(`a`,`
 text-decoration: none;
 color: inherit;
 outline: none;
 `,[S(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),y(`dropdown-option-body`,`
 display: flex;
 cursor: pointer;
 position: relative;
 height: var(--n-option-height);
 line-height: var(--n-option-height);
 font-size: var(--n-font-size);
 color: var(--n-option-text-color);
 transition: color .3s var(--n-bezier);
 `,[S(`&::before`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 left: 4px;
 right: 4px;
 transition: background-color .3s var(--n-bezier);
 border-radius: var(--n-border-radius);
 `),_(`disabled`,[N(`pending`,`
 color: var(--n-option-text-color-hover);
 `,[j(`prefix, suffix`,`
 color: var(--n-option-text-color-hover);
 `),S(`&::before`,`background-color: var(--n-option-color-hover);`)]),N(`active`,`
 color: var(--n-option-text-color-active);
 `,[j(`prefix, suffix`,`
 color: var(--n-option-text-color-active);
 `),S(`&::before`,`background-color: var(--n-option-color-active);`)]),N(`child-active`,`
 color: var(--n-option-text-color-child-active);
 `,[j(`prefix, suffix`,`
 color: var(--n-option-text-color-child-active);
 `)])]),N(`disabled`,`
 cursor: not-allowed;
 opacity: var(--n-option-opacity-disabled);
 `),N(`group`,`
 font-size: calc(var(--n-font-size) - 1px);
 color: var(--n-group-header-text-color);
 `,[j(`prefix`,`
 width: calc(var(--n-option-prefix-width) / 2);
 `,[N(`show-icon`,`
 width: calc(var(--n-option-icon-prefix-width) / 2);
 `)])]),j(`prefix`,`
 width: var(--n-option-prefix-width);
 display: flex;
 justify-content: center;
 align-items: center;
 color: var(--n-prefix-color);
 transition: color .3s var(--n-bezier);
 z-index: 1;
 `,[N(`show-icon`,`
 width: var(--n-option-icon-prefix-width);
 `),y(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),j(`label`,`
 white-space: nowrap;
 flex: 1;
 z-index: 1;
 `),j(`suffix`,`
 box-sizing: border-box;
 flex-grow: 0;
 flex-shrink: 0;
 display: flex;
 justify-content: flex-end;
 align-items: center;
 min-width: var(--n-option-suffix-width);
 padding: 0 8px;
 transition: color .3s var(--n-bezier);
 color: var(--n-suffix-color);
 z-index: 1;
 `,[N(`has-submenu`,`
 width: var(--n-option-icon-suffix-width);
 `),y(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),y(`dropdown-menu`,`pointer-events: all;`)]),y(`dropdown-offset-container`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: -4px;
 bottom: -4px;
 `)]),y(`dropdown-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 4px 0;
 `),y(`dropdown-menu-wrapper`,`
 transform-origin: var(--v-transform-origin);
 width: fit-content;
 `),S(`>`,[y(`scrollbar`,`
 height: inherit;
 max-height: inherit;
 `)]),_(`scrollable`,`
 padding: var(--n-padding);
 `),N(`scrollable`,[j(`content`,`
 padding: var(--n-padding);
 `)])]),Oe={animated:{type:Boolean,default:!0},keyboard:{type:Boolean,default:!0},size:String,inverted:Boolean,placement:{type:String,default:`bottom`},onSelect:[Function,Array],options:{type:Array,default:()=>[]},menuProps:Function,showArrow:Boolean,renderLabel:Function,renderIcon:Function,renderOption:Function,nodeProps:Function,labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},value:[String,Number]},ke=Object.keys(R),Ae={...R,...Oe,...w.props},je=h({name:`Dropdown`,inheritAttrs:!1,props:Ae,setup(e){let t=x(!1),n=oe(A(e,`show`),t),r=a(()=>{let{keyField:t,childrenField:n}=e;return te(e.options,{getKey(e){return e[t]},getDisabled(e){return e.disabled===!0},getIgnored(e){return e.type===`divider`||e.type===`render`},getChildren(e){return e[n]}})}),i=a(()=>r.value.treeNodes),o=x(null),s=x(null),c=x(null),l=a(()=>o.value??s.value??c.value??null),u=a(()=>r.value.getPath(l.value).keyPath),f=a(()=>r.value.getPath(e.value).keyPath),p=v(()=>e.keyboard&&n.value);me({keydown:{ArrowUp:{prevent:!0,handler:N},ArrowRight:{prevent:!0,handler:j},ArrowDown:{prevent:!0,handler:P},ArrowLeft:{prevent:!0,handler:D},Enter:{prevent:!0,handler:F},Escape:E}},p);let{mergedClsPrefixRef:m,inlineThemeDisabled:h,mergedComponentPropsRef:g}=d(e),_=a(()=>e.size||g?.value?.Dropdown?.size||`medium`),y=w(`Dropdown`,`-dropdown`,De,pe,e,m);O(J,{labelFieldRef:A(e,`labelField`),childrenFieldRef:A(e,`childrenField`),renderLabelRef:A(e,`renderLabel`),renderIconRef:A(e,`renderIcon`),hoverKeyRef:o,keyboardKeyRef:s,lastToggledSubmenuKeyRef:c,pendingKeyPathRef:u,activeKeyPathRef:f,animatedRef:A(e,`animated`),mergedShowRef:n,nodePropsRef:A(e,`nodeProps`),renderOptionRef:A(e,`renderOption`),menuPropsRef:A(e,`menuProps`),doSelect:S,doUpdateShow:C}),k(n,t=>{!e.animated&&!t&&T()});function S(t,n){let{onSelect:r}=e;r&&V(r,t,n)}function C(n){let{"onUpdate:show":r,onUpdateShow:i}=e;r&&V(r,n),i&&V(i,n),t.value=n}function T(){o.value=null,s.value=null,c.value=null}function E(){C(!1)}function D(){L(`left`)}function j(){L(`right`)}function N(){L(`up`)}function P(){L(`down`)}function F(){let e=I();e?.isLeaf&&n.value&&(S(e.key,e.rawNode),C(!1))}function I(){let{value:e}=r,{value:t}=l;return!e||t===null?null:e.getNode(t)??null}function L(e){let{value:t}=l,{value:{getFirstAvailableNode:n}}=r,i=null;if(t===null){let e=n();e!==null&&(i=e.key)}else{let t=I();if(t){let n;switch(e){case`down`:n=t.getNext();break;case`up`:n=t.getPrev();break;case`right`:n=t.getChild();break;case`left`:n=t.getParent()}n&&(i=n.key)}}i!==null&&(o.value=null,s.value=i)}let R=a(()=>{let{inverted:t}=e,n=_.value,{common:{cubicBezierEaseInOut:r},self:i}=y.value,{padding:a,dividerColor:o,borderRadius:s,optionOpacityDisabled:c,[M(`optionIconSuffixWidth`,n)]:l,[M(`optionSuffixWidth`,n)]:u,[M(`optionIconPrefixWidth`,n)]:d,[M(`optionPrefixWidth`,n)]:f,[M(`fontSize`,n)]:p,[M(`optionHeight`,n)]:m,[M(`optionIconSize`,n)]:h}=i,g={"--n-bezier":r,"--n-font-size":p,"--n-padding":a,"--n-border-radius":s,"--n-option-height":m,"--n-option-prefix-width":f,"--n-option-icon-prefix-width":d,"--n-option-suffix-width":u,"--n-option-icon-suffix-width":l,"--n-option-icon-size":h,"--n-divider-color":o,"--n-option-opacity-disabled":c};return t?(g[`--n-color`]=i.colorInverted,g[`--n-option-color-hover`]=i.optionColorHoverInverted,g[`--n-option-color-active`]=i.optionColorActiveInverted,g[`--n-option-text-color`]=i.optionTextColorInverted,g[`--n-option-text-color-hover`]=i.optionTextColorHoverInverted,g[`--n-option-text-color-active`]=i.optionTextColorActiveInverted,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActiveInverted,g[`--n-prefix-color`]=i.prefixColorInverted,g[`--n-suffix-color`]=i.suffixColorInverted,g[`--n-group-header-text-color`]=i.groupHeaderTextColorInverted):(g[`--n-color`]=i.color,g[`--n-option-color-hover`]=i.optionColorHover,g[`--n-option-color-active`]=i.optionColorActive,g[`--n-option-text-color`]=i.optionTextColor,g[`--n-option-text-color-hover`]=i.optionTextColorHover,g[`--n-option-text-color-active`]=i.optionTextColorActive,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActive,g[`--n-prefix-color`]=i.prefixColor,g[`--n-suffix-color`]=i.suffixColor,g[`--n-group-header-text-color`]=i.groupHeaderTextColor),g}),z=h?b(`dropdown`,a(()=>`${_.value[0]}${e.inverted?`i`:``}`),R,e):void 0;return{mergedClsPrefix:m,mergedTheme:y,mergedSize:_,tmNodes:i,mergedShow:n,handleAfterLeave:()=>{e.animated&&T()},doUpdateShow:C,cssVars:h?void 0:R,themeClass:z?.themeClass,onRender:z?.onRender}},render(){let e=(e,t,n,r,a)=>{let{mergedClsPrefix:s,menuProps:c}=this;this.onRender?.();let l=c?.(void 0,this.tmNodes.map(e=>e.rawNode))||{},u={ref:K(t),class:[e,`${s}-dropdown`,`${s}-dropdown--${this.mergedSize}-size`,this.themeClass],clsPrefix:s,tmNodes:this.tmNodes,style:[...n,this.cssVars],showArrow:this.showArrow,arrowStyle:this.arrowStyle,scrollable:this.scrollable,onMouseenter:r,onMouseleave:a};return i(Ee,o(this.$attrs,u,l))},{mergedTheme:n}=this,r={show:this.mergedShow,theme:n.peers.Popover,themeOverrides:n.peerOverrides.Popover,internalOnAfterLeave:this.handleAfterLeave,internalRenderBody:e,onUpdateShow:this.doUpdateShow,"onUpdate:show":void 0};return E(),p(z,ae(this.$props,ke,r),{_:1,trigger:t(()=>this.$slots.default?.())},16)}});export{G as i,ge as n,K as r,je as t};