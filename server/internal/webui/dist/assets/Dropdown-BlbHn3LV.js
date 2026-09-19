import{$t as e,At as t,Bt as n,D as r,Ft as i,Jt as a,Lt as o,Mt as s,Ot as c,Rt as l,Ut as u,Wt as d,Yt as f,_ as p,an as m,at as h,ct as g,dt as _,g as v,gt as y,h as b,k as x,kt as S,lt as C,o as w,on as T,ot as E,s as D,sn as O,un as k,ut as A,w as j,wt as M,y as N}from"./vue-i18n-CgKceEWx.js";import{a as P,d as F,g as I,i as L,r as R,s as ee,t as te,u as ne}from"./create-Crm1UAxV.js";import{t as re}from"./format-length-04qCp4Ww.js";import{g as z,h as B,m as V,r as ie}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{c as ae,f as oe,g as se,h as H,p as ce,t as le,y as ue}from"./fade-in-scale-up.cssr-7akLf3Ad.js";import{r as U}from"./FadeInExpandTransition-D8IUT6y9.js";import{D as W,_ as de,p as fe,v as pe}from"./index-kjlF3dpY.js";function me(t={},n){let r=m({ctrl:!1,command:!1,win:!1,shift:!1,tab:!1}),{keydown:i,keyup:a}=t,o=e=>{switch(e.key){case`Control`:r.ctrl=!0;break;case`Meta`:r.command=!0,r.win=!0;break;case`Shift`:r.shift=!0;break;case`Tab`:r.tab=!0}i!==void 0&&Object.keys(i).forEach(t=>{if(t!==e.key)return;let n=i[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},s=e=>{switch(e.key){case`Control`:r.ctrl=!1;break;case`Meta`:r.command=!1,r.win=!1;break;case`Shift`:r.shift=!1;break;case`Tab`:r.tab=!1}a!==void 0&&Object.keys(a).forEach(t=>{if(t!==e.key)return;let n=a[t];if(typeof n==`function`)n(e);else{let{stop:t=!1,prevent:r=!1}=n;t&&e.stopPropagation(),r&&e.preventDefault(),n.handler(e)}})},c=()=>{(n===void 0||n.value)&&(z(`keydown`,document,o),z(`keyup`,document,s)),n!==void 0&&e(n,e=>{e?(z(`keydown`,document,o),z(`keyup`,document,s)):(B(`keydown`,document,o),B(`keyup`,document,s))})};return ce()?(u(c),d(()=>{(n===void 0||n.value)&&(B(`keydown`,document,o),B(`keyup`,document,s))})):c(),T(r)}var G=i({name:`ChevronRight`,render(){return(()=>{let e=b(`6ab04425f4fcb756`);return e[0]||=S(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[S(`path`,{d:`M5.64645 3.14645C5.45118 3.34171 5.45118 3.65829 5.64645 3.85355L9.79289 8L5.64645 12.1464C5.45118 12.3417 5.45118 12.6583 5.64645 12.8536C5.84171 13.0488 6.15829 13.0488 6.35355 12.8536L10.8536 8.35355C11.0488 8.15829 11.0488 7.84171 10.8536 7.64645L6.35355 3.14645C6.15829 2.95118 5.84171 2.95118 5.64645 3.14645Z`,fill:`currentColor`})],-1)})()}});function K(e){return t=>{e.value=t?t.$el:null}}var he={...L,...w.props},ge=i({name:`Tooltip`,props:he,slots:Object,__popover__:!0,setup(e){let{mergedClsPrefixRef:t}=j(e),n=w(`Tooltip`,`-tooltip`,void 0,de,e,t),r=O(null);return{syncPosition(){r.value.syncPosition()},setShow(e){r.value.setShow(e)},popoverRef:r,mergedTheme:n,popoverThemeOverrides:c(()=>n.value.self)}},render(){let{mergedTheme:e,internalExtraClass:t}=this;return o(R,{...this.$props,theme:e.peers.Popover,themeOverrides:e.peerOverrides.Popover,builtinThemeOverrides:this.popoverThemeOverrides,internalExtraClass:t.concat(`tooltip`),ref:`popoverRef`},this.$slots)}}),_e=E(`icon`,`
 height: 1em;
 width: 1em;
 line-height: 1em;
 text-align: center;
 display: inline-block;
 position: relative;
 fill: currentColor;
`,[C(`color-transition`,{transition:`color .3s var(--n-bezier)`}),C(`depth`,{color:`var(--n-color)`},[h(`svg`,{opacity:`var(--n-opacity)`,transition:`opacity .3s var(--n-bezier)`})]),h(`svg`,{height:`1em`,width:`1em`})]),ve={...w.props,depth:[String,Number],size:[Number,String],color:String,component:[Object,Function]},ye=i({_n_icon__:!0,name:`Icon`,inheritAttrs:!1,props:ve,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=j(e),r=w(`Icon`,`-icon`,_e,fe,e,t),i=c(()=>{let{depth:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value;if(t!==void 0){let{color:e,[`opacity${t}Depth`]:r}=i;return{"--n-bezier":n,"--n-color":e,"--n-opacity":r}}return{"--n-bezier":n,"--n-color":``,"--n-opacity":``}}),a=n?D(`icon`,c(()=>`${e.depth||`d`}`),i,e):void 0;return{mergedClsPrefix:t,mergedStyle:c(()=>{let{size:t,color:n}=e;return{fontSize:re(t),color:n}}),cssVars:n?void 0:i,themeClass:a?.themeClass,onRender:a?.onRender}},render(){let{$parent:e,depth:t,mergedClsPrefix:r,component:i,onRender:a,themeClass:s}=this;return e?.$options?._n_icon__&&x(`icon`,"don't wrap `n-icon` inside `n-icon`"),a?.(),o(`i`,n(this.$attrs,{role:`img`,class:[`${r}-icon`,s,{[`${r}-icon--depth`]:t,[`${r}-icon--color-transition`]:t!==void 0}],style:[this.cssVars,this.mergedStyle]}),i?o(i):this.$slots.default?.())}}),q=r(`n-dropdown-menu`),J=r(`n-dropdown`),Y=r(`n-dropdown-option`),X=i({name:`DropdownDivider`,props:{clsPrefix:{type:String,required:!0}},render(){return a(),s(`div`,{class:v(`${this.clsPrefix}-dropdown-divider`)},null,2)}});function Z(e,t){return e.type===`submenu`||e.type===void 0&&e[t]!==void 0}function be(e){return e.type===`group`}function Q(e){return e.type===`divider`}function xe(e){return e.type===`render`}function Se(t,n,r){if(!n)return t;let i=O(t.value),a=null;return e(t,e=>{a!==null&&window.clearTimeout(a),e===!0?r&&!r.value?i.value=!0:a=window.setTimeout(()=>{i.value=!0},n):i.value=!1}),i}var $=i({name:`DropdownOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null},placement:{type:String,default:`right-start`},props:Object,scrollable:Boolean},setup(e){let t=l(J),{hoverKeyRef:n,keyboardKeyRef:r,lastToggledSubmenuKeyRef:i,pendingKeyPathRef:a,activeKeyPathRef:o,animatedRef:s,mergedShowRef:u,renderLabelRef:d,renderIconRef:p,labelFieldRef:m,childrenFieldRef:h,renderOptionRef:g,nodePropsRef:_,menuPropsRef:v}=t,y=l(Y,null),b=l(q),x=l(H),S=c(()=>e.tmNode.rawNode),C=c(()=>{let{value:t}=h;return Z(e.tmNode.rawNode,t)}),w=c(()=>{let{disabled:t}=e.tmNode;return t}),T=Se(c(()=>{if(!C.value)return!1;let{key:t,disabled:o}=e.tmNode;if(o)return!1;let{value:s}=n,{value:c}=r,{value:l}=i,{value:u}=a;return s===null?c===null?l!==null&&u.includes(t):u.includes(t)&&u[u.length-1]!==t:u.includes(t)}),300,c(()=>r.value===null&&!s.value)),E=c(()=>!!y?.enteringSubmenuRef.value),D=O(!1);f(Y,{enteringSubmenuRef:D});function k(){D.value=!0}function A(){D.value=!1}function j(){let{parentKey:t,tmNode:a}=e;a.disabled||u.value&&(i.value=t,r.value=null,n.value=a.key)}function M(){let{tmNode:t}=e;t.disabled||u.value&&n.value!==t.key&&j()}function N(t){if(e.tmNode.disabled||!u.value)return;let{relatedTarget:r}=t;r&&!I({target:r},`dropdownOption`)&&!I({target:r},`scrollbarRail`)&&(n.value=null)}function P(){let{value:n}=C,{tmNode:r}=e;u.value&&!n&&!r.disabled&&(t.doSelect(r.key,r.rawNode),t.doUpdateShow(!1))}return{labelField:m,renderLabel:d,renderIcon:p,siblingHasIcon:b.showIconRef,siblingHasSubmenu:b.hasSubmenuRef,menuProps:v,popoverBody:x,animated:s,mergedShowSubmenu:c(()=>T.value&&!E.value),rawNode:S,hasSubmenu:C,pending:U(()=>{let{value:t}=a,{key:n}=e.tmNode;return t.includes(n)}),childActive:U(()=>{let{value:t}=o,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r!==-1&&r<t.length-1}),active:U(()=>{let{value:t}=o,{key:n}=e.tmNode,r=t.findIndex(e=>n===e);return r!==-1&&r===t.length-1}),mergedDisabled:w,renderOption:g,nodeProps:_,handleClick:P,handleMouseMove:M,handleMouseEnter:j,handleMouseLeave:N,handleSubmenuBeforeEnter:k,handleSubmenuAfterEnter:A}},render(){let{animated:e,rawNode:r,mergedShowSubmenu:i,clsPrefix:c,siblingHasIcon:l,siblingHasSubmenu:u,renderLabel:d,renderIcon:f,renderOption:m,nodeProps:h,props:g,scrollable:_}=this,b=null;if(i){let e=this.menuProps?.(r,r.children);b=(r=>(a(),t(Ee,n({key:1},e,{clsPrefix:c,scrollable:this.scrollable,tmNodes:this.tmNode.children,parentKey:this.tmNode.key}),null,16,[`clsPrefix`,`scrollable`,`tmNodes`,`parentKey`])))(b)}let x={class:[`${c}-dropdown-option-body`,this.pending&&`${c}-dropdown-option-body--pending`,this.active&&`${c}-dropdown-option-body--active`,this.childActive&&`${c}-dropdown-option-body--child-active`,this.mergedDisabled&&`${c}-dropdown-option-body--disabled`],onMousemove:this.handleMouseMove,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onClick:this.handleClick},S=h?.(r),C=(a(),s(`div`,n({class:[`${c}-dropdown-option`,S?.class],"data-dropdown-option":!0},S),[N(()=>o(`div`,n(x,g),[(a(),s(`div`,{class:v([`${c}-dropdown-option-body__prefix`,l&&`${c}-dropdown-option-body__prefix--show-icon`])},[N(()=>[f?f(r):W(r.icon)])],2)),(a(),s(`div`,{"data-dropdown-option":!0,class:v(`${c}-dropdown-option-body__label`)},[d?(a(),s(M,{key:0},[N(()=>d(r))],64)):(a(),s(M,{key:1},[N(()=>W(r[this.labelField]??r.title))],64))],2)),(a(),s(`div`,{"data-dropdown-option":!0,class:v([`${c}-dropdown-option-body__suffix`,u&&`${c}-dropdown-option-body__suffix--has-submenu`])},[this.hasSubmenu?(a(),t(ye,{key:0},{_:1,default:p(()=>(a(),t(G)))})):N(()=>null)],2))])),this.hasSubmenu?(a(),t(F,{key:0},{default:()=>[(a(),t(ne,null,{default:()=>(a(),s(`div`,{class:v(`${c}-dropdown-offset-container`)},[(a(),t(ee,{show:this.mergedShowSubmenu,placement:this.placement,to:_&&this.popoverBody||void 0,teleportDisabled:!_},{default:()=>(a(),s(`div`,{class:v(`${c}-dropdown-menu-wrapper`)},[e?(a(),t(y,{key:0,onBeforeEnter:this.handleSubmenuBeforeEnter,onAfterEnter:this.handleSubmenuAfterEnter,name:`fade-in-scale-up-transition`,appear:!0},{default:()=>b},1032,[`onBeforeEnter`,`onAfterEnter`])):(a(),s(M,{key:1},[N(()=>b)],64))],2))},1032,[`show`,`placement`,`to`,`teleportDisabled`]))],2))},1024))]},1024)):N(()=>null)],16));return m?m({node:C,option:r}):C}}),Ce=i({name:`DropdownGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{showIconRef:e,hasSubmenuRef:t}=l(q),{renderLabelRef:n,labelFieldRef:r,nodePropsRef:i,renderOptionRef:a}=l(J);return{labelField:r,showIcon:e,hasSubmenu:t,renderLabel:n,nodeProps:i,renderOption:a}},render(){let{clsPrefix:e,hasSubmenu:t,showIcon:r,nodeProps:i,renderLabel:o,renderOption:c}=this,{rawNode:l}=this.tmNode,u=(a(),s(`div`,n({class:`${e}-dropdown-option`},i?.(l)),[S(`div`,{class:v(`${e}-dropdown-option-body ${e}-dropdown-option-body--group`)},[S(`div`,{"data-dropdown-option":!0,class:v([`${e}-dropdown-option-body__prefix`,r&&`${e}-dropdown-option-body__prefix--show-icon`])},[N(()=>W(l.icon))],2),S(`div`,{class:v(`${e}-dropdown-option-body__label`),"data-dropdown-option":!0},[o?(a(),s(M,{key:0},[N(()=>o(l))],64)):(a(),s(M,{key:1},[N(()=>W(l.title??l[this.labelField]))],64))],2),S(`div`,{class:v([`${e}-dropdown-option-body__suffix`,t&&`${e}-dropdown-option-body__suffix--has-submenu`]),"data-dropdown-option":!0},null,2)],2)],16));return c?c({node:u,option:l}):u}}),we=i({name:`NDropdownGroup`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0},parentKey:{type:[String,Number],default:null}},render(){let{tmNode:e,parentKey:n,clsPrefix:r}=this,{children:i}=e;return a(),s(M,null,[(a(),t(Ce,{clsPrefix:r,tmNode:e,key:e.key},null,8,[`clsPrefix`,`tmNode`])),N(()=>i?.map(e=>{let{rawNode:i}=e;return i.show===!1?null:Q(i)?o(X,{clsPrefix:r,key:e.key}):e.isGroup?(x(`dropdown`,"`group` node is not allowed to be put in `group` node."),null):(a(),t($,{clsPrefix:r,tmNode:e,parentKey:n,key:e.key},null,8,[`clsPrefix`,`tmNode`,`parentKey`]))}))],64)}}),Te=i({name:`DropdownRenderOption`,props:{tmNode:{type:Object,required:!0}},render(){let{rawNode:{render:e,props:t}}=this.tmNode;return o(`div`,t,[e?.()])}}),Ee=i({name:`DropdownMenu`,props:{scrollable:Boolean,showArrow:Boolean,arrowStyle:[String,Object],clsPrefix:{type:String,required:!0},tmNodes:{type:Array,default:()=>[]},parentKey:{type:[String,Number],default:null}},setup(e){let{renderIconRef:t,childrenFieldRef:n}=l(J);f(q,{showIconRef:c(()=>{let n=t.value;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>n?n(e):e.icon);let{rawNode:t}=e;return n?n(t):t.icon})}),hasSubmenuRef:c(()=>{let{value:t}=n;return e.tmNodes.some(e=>{if(e.isGroup)return e.children?.some(({rawNode:e})=>Z(e,t));let{rawNode:n}=e;return Z(n,t)})})});let r=O(null);return f(se,null),f(ue,null),f(H,r),{bodyRef:r}},render(){let{parentKey:e,clsPrefix:n,scrollable:r}=this,i=this.tmNodes.map(i=>{let{rawNode:o}=i;return o.show===!1?null:xe(o)?(a(),t(Te,{tmNode:i,key:i.key},null,8,[`tmNode`])):Q(o)?(a(),t(X,{clsPrefix:n,key:i.key},null,8,[`clsPrefix`])):be(o)?(a(),t(we,{clsPrefix:n,tmNode:i,parentKey:e,key:i.key},null,8,[`clsPrefix`,`tmNode`,`parentKey`])):(a(),t($,{clsPrefix:n,tmNode:i,parentKey:e,key:i.key,props:o.props,scrollable:r},null,8,[`clsPrefix`,`tmNode`,`parentKey`,`props`,`scrollable`]))});return a(),s(`div`,{class:v([`${n}-dropdown-menu`,r&&`${n}-dropdown-menu--scrollable`]),ref:`bodyRef`},[r?(a(),t(ie,{key:0,contentClass:`${n}-dropdown-menu__content`},{default:()=>i},1032,[`contentClass`])):(a(),s(M,{key:1},[N(()=>i)],64)),this.showArrow?(a(),s(M,{key:2},[N(()=>P({clsPrefix:n,arrowStyle:this.arrowStyle,arrowClass:void 0,arrowWrapperClass:void 0,arrowWrapperStyle:void 0}))],64)):N(()=>null)],2)}}),De=E(`dropdown-menu`,`
 transform-origin: var(--v-transform-origin);
 background-color: var(--n-color);
 border-radius: var(--n-border-radius);
 box-shadow: var(--n-box-shadow);
 position: relative;
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
`,[le(),E(`dropdown-option`,`
 position: relative;
 `,[h(`a`,`
 text-decoration: none;
 color: inherit;
 outline: none;
 `,[h(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),E(`dropdown-option-body`,`
 display: flex;
 cursor: pointer;
 position: relative;
 height: var(--n-option-height);
 line-height: var(--n-option-height);
 font-size: var(--n-font-size);
 color: var(--n-option-text-color);
 transition: color .3s var(--n-bezier);
 `,[h(`&::before`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 left: 4px;
 right: 4px;
 transition: background-color .3s var(--n-bezier);
 border-radius: var(--n-border-radius);
 `),A(`disabled`,[C(`pending`,`
 color: var(--n-option-text-color-hover);
 `,[g(`prefix, suffix`,`
 color: var(--n-option-text-color-hover);
 `),h(`&::before`,`background-color: var(--n-option-color-hover);`)]),C(`active`,`
 color: var(--n-option-text-color-active);
 `,[g(`prefix, suffix`,`
 color: var(--n-option-text-color-active);
 `),h(`&::before`,`background-color: var(--n-option-color-active);`)]),C(`child-active`,`
 color: var(--n-option-text-color-child-active);
 `,[g(`prefix, suffix`,`
 color: var(--n-option-text-color-child-active);
 `)])]),C(`disabled`,`
 cursor: not-allowed;
 opacity: var(--n-option-opacity-disabled);
 `),C(`group`,`
 font-size: calc(var(--n-font-size) - 1px);
 color: var(--n-group-header-text-color);
 `,[g(`prefix`,`
 width: calc(var(--n-option-prefix-width) / 2);
 `,[C(`show-icon`,`
 width: calc(var(--n-option-icon-prefix-width) / 2);
 `)])]),g(`prefix`,`
 width: var(--n-option-prefix-width);
 display: flex;
 justify-content: center;
 align-items: center;
 color: var(--n-prefix-color);
 transition: color .3s var(--n-bezier);
 z-index: 1;
 `,[C(`show-icon`,`
 width: var(--n-option-icon-prefix-width);
 `),E(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),g(`label`,`
 white-space: nowrap;
 flex: 1;
 z-index: 1;
 `),g(`suffix`,`
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
 `,[C(`has-submenu`,`
 width: var(--n-option-icon-suffix-width);
 `),E(`icon`,`
 font-size: var(--n-option-icon-size);
 `)]),E(`dropdown-menu`,`pointer-events: all;`)]),E(`dropdown-offset-container`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: -4px;
 bottom: -4px;
 `)]),E(`dropdown-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 4px 0;
 `),E(`dropdown-menu-wrapper`,`
 transform-origin: var(--v-transform-origin);
 width: fit-content;
 `),h(`>`,[E(`scrollbar`,`
 height: inherit;
 max-height: inherit;
 `)]),A(`scrollable`,`
 padding: var(--n-padding);
 `),C(`scrollable`,[g(`content`,`
 padding: var(--n-padding);
 `)])]),Oe={animated:{type:Boolean,default:!0},keyboard:{type:Boolean,default:!0},size:String,inverted:Boolean,placement:{type:String,default:`bottom`},onSelect:[Function,Array],options:{type:Array,default:()=>[]},menuProps:Function,showArrow:Boolean,renderLabel:Function,renderIcon:Function,renderOption:Function,nodeProps:Function,labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},value:[String,Number]},ke=Object.keys(L),Ae={...L,...Oe,...w.props},je=i({name:`Dropdown`,inheritAttrs:!1,props:Ae,setup(t){let n=O(!1),r=oe(k(t,`show`),n),i=c(()=>{let{keyField:e,childrenField:n}=t;return te(t.options,{getKey(t){return t[e]},getDisabled(e){return e.disabled===!0},getIgnored(e){return e.type===`divider`||e.type===`render`},getChildren(e){return e[n]}})}),a=c(()=>i.value.treeNodes),o=O(null),s=O(null),l=O(null),u=c(()=>o.value??s.value??l.value??null),d=c(()=>i.value.getPath(u.value).keyPath),p=c(()=>i.value.getPath(t.value).keyPath),m=U(()=>t.keyboard&&r.value);me({keydown:{ArrowUp:{prevent:!0,handler:M},ArrowRight:{prevent:!0,handler:A},ArrowDown:{prevent:!0,handler:N},ArrowLeft:{prevent:!0,handler:E},Enter:{prevent:!0,handler:P},Escape:T}},m);let{mergedClsPrefixRef:h,inlineThemeDisabled:g,mergedComponentPropsRef:v}=j(t),y=c(()=>t.size||v?.value?.Dropdown?.size||`medium`),b=w(`Dropdown`,`-dropdown`,De,pe,t,h);f(J,{labelFieldRef:k(t,`labelField`),childrenFieldRef:k(t,`childrenField`),renderLabelRef:k(t,`renderLabel`),renderIconRef:k(t,`renderIcon`),hoverKeyRef:o,keyboardKeyRef:s,lastToggledSubmenuKeyRef:l,pendingKeyPathRef:d,activeKeyPathRef:p,animatedRef:k(t,`animated`),mergedShowRef:r,nodePropsRef:k(t,`nodeProps`),renderOptionRef:k(t,`renderOption`),menuPropsRef:k(t,`menuProps`),doSelect:x,doUpdateShow:S}),e(r,e=>{!t.animated&&!e&&C()});function x(e,n){let{onSelect:r}=t;r&&V(r,e,n)}function S(e){let{"onUpdate:show":r,onUpdateShow:i}=t;r&&V(r,e),i&&V(i,e),n.value=e}function C(){o.value=null,s.value=null,l.value=null}function T(){S(!1)}function E(){I(`left`)}function A(){I(`right`)}function M(){I(`up`)}function N(){I(`down`)}function P(){let e=F();e?.isLeaf&&r.value&&(x(e.key,e.rawNode),S(!1))}function F(){let{value:e}=i,{value:t}=u;return!e||t===null?null:e.getNode(t)??null}function I(e){let{value:t}=u,{value:{getFirstAvailableNode:n}}=i,r=null;if(t===null){let e=n();e!==null&&(r=e.key)}else{let t=F();if(t){let n;switch(e){case`down`:n=t.getNext();break;case`up`:n=t.getPrev();break;case`right`:n=t.getChild();break;case`left`:n=t.getParent()}n&&(r=n.key)}}r!==null&&(o.value=null,s.value=r)}let L=c(()=>{let{inverted:e}=t,n=y.value,{common:{cubicBezierEaseInOut:r},self:i}=b.value,{padding:a,dividerColor:o,borderRadius:s,optionOpacityDisabled:c,[_(`optionIconSuffixWidth`,n)]:l,[_(`optionSuffixWidth`,n)]:u,[_(`optionIconPrefixWidth`,n)]:d,[_(`optionPrefixWidth`,n)]:f,[_(`fontSize`,n)]:p,[_(`optionHeight`,n)]:m,[_(`optionIconSize`,n)]:h}=i,g={"--n-bezier":r,"--n-font-size":p,"--n-padding":a,"--n-border-radius":s,"--n-option-height":m,"--n-option-prefix-width":f,"--n-option-icon-prefix-width":d,"--n-option-suffix-width":u,"--n-option-icon-suffix-width":l,"--n-option-icon-size":h,"--n-divider-color":o,"--n-option-opacity-disabled":c};return e?(g[`--n-color`]=i.colorInverted,g[`--n-option-color-hover`]=i.optionColorHoverInverted,g[`--n-option-color-active`]=i.optionColorActiveInverted,g[`--n-option-text-color`]=i.optionTextColorInverted,g[`--n-option-text-color-hover`]=i.optionTextColorHoverInverted,g[`--n-option-text-color-active`]=i.optionTextColorActiveInverted,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActiveInverted,g[`--n-prefix-color`]=i.prefixColorInverted,g[`--n-suffix-color`]=i.suffixColorInverted,g[`--n-group-header-text-color`]=i.groupHeaderTextColorInverted):(g[`--n-color`]=i.color,g[`--n-option-color-hover`]=i.optionColorHover,g[`--n-option-color-active`]=i.optionColorActive,g[`--n-option-text-color`]=i.optionTextColor,g[`--n-option-text-color-hover`]=i.optionTextColorHover,g[`--n-option-text-color-active`]=i.optionTextColorActive,g[`--n-option-text-color-child-active`]=i.optionTextColorChildActive,g[`--n-prefix-color`]=i.prefixColor,g[`--n-suffix-color`]=i.suffixColor,g[`--n-group-header-text-color`]=i.groupHeaderTextColor),g}),R=g?D(`dropdown`,c(()=>`${y.value[0]}${t.inverted?`i`:``}`),L,t):void 0;return{mergedClsPrefix:h,mergedTheme:b,mergedSize:y,tmNodes:a,mergedShow:r,handleAfterLeave:()=>{t.animated&&C()},doUpdateShow:S,cssVars:g?void 0:L,themeClass:R?.themeClass,onRender:R?.onRender}},render(){let e=(e,t,r,i,a)=>{let{mergedClsPrefix:s,menuProps:c}=this;this.onRender?.();let l=c?.(void 0,this.tmNodes.map(e=>e.rawNode))||{},u={ref:K(t),class:[e,`${s}-dropdown`,`${s}-dropdown--${this.mergedSize}-size`,this.themeClass],clsPrefix:s,tmNodes:this.tmNodes,style:[...r,this.cssVars],showArrow:this.showArrow,arrowStyle:this.arrowStyle,scrollable:this.scrollable,onMouseenter:i,onMouseleave:a};return o(Ee,n(this.$attrs,u,l))},{mergedTheme:r}=this,i={show:this.mergedShow,theme:r.peers.Popover,themeOverrides:r.peerOverrides.Popover,internalOnAfterLeave:this.handleAfterLeave,internalRenderBody:e,onUpdateShow:this.doUpdateShow,"onUpdate:show":void 0};return a(),t(R,ae(this.$props,ke,i),{_:1,trigger:p(()=>this.$slots.default?.())},16)}});export{G as i,ge as n,K as r,je as t};