import{At as e,Bt as t,D as n,Ft as r,Jt as i,Lt as a,Mt as o,Ot as s,Rt as c,Yt as l,at as u,ct as d,en as f,g as p,h as m,i as h,kt as g,lt as _,o as v,ot as y,pn as b,s as x,sn as S,un as ee,ut as C,w as te,wt as w,y as T}from"./vue-i18n-CgKceEWx.js";import{o as E,p as ne,t as D}from"./create-Crm1UAxV.js";import{S as O,i as k,m as A}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{c as j,f as re}from"./fade-in-scale-up.cssr-7akLf3Ad.js";import{i as ie,r as M,t as N}from"./FadeInExpandTransition-D8IUT6y9.js";import{n as P,t as F}from"./Dropdown-BlbHn3LV.js";import{D as I,M as ae,o as oe}from"./index-kjlF3dpY.js";var se=n(`n-layout-sider`),L={type:String,default:`static`},R=n(`n-menu`),z=n(`n-submenu`),B=n(`n-menu-item-group`),V=[u(`&::before`,`background-color: var(--n-item-color-hover);`),d(`arrow`,`
 color: var(--n-arrow-color-hover);
 `),d(`icon`,`
 color: var(--n-item-icon-color-hover);
 `),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover);
 `,[u(`a`,`
 color: var(--n-item-text-color-hover);
 `),d(`extra`,`
 color: var(--n-item-text-color-hover);
 `)])],H=[d(`icon`,`
 color: var(--n-item-icon-color-hover-horizontal);
 `),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover-horizontal);
 `,[u(`a`,`
 color: var(--n-item-text-color-hover-horizontal);
 `),d(`extra`,`
 color: var(--n-item-text-color-hover-horizontal);
 `)])],ce=u([y(`menu`,`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[_(`horizontal`,`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[y(`submenu`,`margin: 0;`),y(`menu-item`,`margin: 0;`),y(`menu-item-content`,`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[u(`&::before`,`display: none;`),_(`selected`,`border-bottom: 2px solid var(--n-border-color-horizontal)`)]),y(`menu-item-content`,[_(`selected`,[d(`icon`,`color: var(--n-item-icon-color-active-horizontal);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-horizontal);
 `,[u(`a`,`color: var(--n-item-text-color-active-horizontal);`),d(`extra`,`color: var(--n-item-text-color-active-horizontal);`)])]),_(`child-active`,`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[y(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[u(`a`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `),d(`extra`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),d(`icon`,`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),C(`disabled`,[C(`selected, child-active`,[u(`&:focus-within`,H)]),_(`selected`,[U(null,[d(`icon`,`color: var(--n-item-icon-color-active-hover-horizontal);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[u(`a`,`color: var(--n-item-text-color-active-hover-horizontal);`),d(`extra`,`color: var(--n-item-text-color-active-hover-horizontal);`)])])]),_(`child-active`,[U(null,[d(`icon`,`color: var(--n-item-icon-color-child-active-hover-horizontal);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[u(`a`,`color: var(--n-item-text-color-child-active-hover-horizontal);`),d(`extra`,`color: var(--n-item-text-color-child-active-hover-horizontal);`)])])]),U(`border-bottom: 2px solid var(--n-border-color-horizontal);`,H)]),y(`menu-item-content-header`,[u(`a`,`color: var(--n-item-text-color-horizontal);`)])])]),C(`responsive`,[y(`menu-item-content-header`,`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),_(`collapsed`,[y(`menu-item-content`,[_(`selected`,[u(`&::before`,`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),y(`menu-item-content-header`,`opacity: 0;`),d(`arrow`,`opacity: 0;`),d(`icon`,`color: var(--n-item-icon-color-collapsed);`)])]),y(`menu-item`,`
 height: var(--n-item-height);
 margin-top: 6px;
 position: relative;
 `),y(`menu-item-content`,`
 box-sizing: border-box;
 line-height: 1.75;
 height: 100%;
 display: grid;
 grid-template-areas: "icon content arrow";
 grid-template-columns: auto 1fr auto;
 align-items: center;
 cursor: pointer;
 position: relative;
 padding-right: 18px;
 transition:
 background-color .3s var(--n-bezier),
 padding-left .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[u(`> *`,`z-index: 1;`),u(`&::before`,`
 z-index: auto;
 content: "";
 background-color: #0000;
 position: absolute;
 left: 8px;
 right: 8px;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),_(`disabled`,`
 opacity: .45;
 cursor: not-allowed;
 `),_(`collapsed`,[d(`arrow`,`transform: rotate(0);`)]),_(`selected`,[u(`&::before`,`background-color: var(--n-item-color-active);`),d(`arrow`,`color: var(--n-arrow-color-active);`),d(`icon`,`color: var(--n-item-icon-color-active);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-active);
 `,[u(`a`,`color: var(--n-item-text-color-active);`),d(`extra`,`color: var(--n-item-text-color-active);`)])]),_(`child-active`,[y(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active);
 `,[u(`a`,`
 color: var(--n-item-text-color-child-active);
 `),d(`extra`,`
 color: var(--n-item-text-color-child-active);
 `)]),d(`arrow`,`
 color: var(--n-arrow-color-child-active);
 `),d(`icon`,`
 color: var(--n-item-icon-color-child-active);
 `)]),C(`disabled`,[C(`selected, child-active`,[u(`&:focus-within`,V)]),_(`selected`,[U(null,[d(`arrow`,`color: var(--n-arrow-color-active-hover);`),d(`icon`,`color: var(--n-item-icon-color-active-hover);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover);
 `,[u(`a`,`color: var(--n-item-text-color-active-hover);`),d(`extra`,`color: var(--n-item-text-color-active-hover);`)])])]),_(`child-active`,[U(null,[d(`arrow`,`color: var(--n-arrow-color-child-active-hover);`),d(`icon`,`color: var(--n-item-icon-color-child-active-hover);`),y(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover);
 `,[u(`a`,`color: var(--n-item-text-color-child-active-hover);`),d(`extra`,`color: var(--n-item-text-color-child-active-hover);`)])])]),_(`selected`,[U(null,[u(`&::before`,`background-color: var(--n-item-color-active-hover);`)])]),U(null,V)]),d(`icon`,`
 grid-area: icon;
 color: var(--n-item-icon-color);
 transition:
 color .3s var(--n-bezier),
 font-size .3s var(--n-bezier),
 margin-right .3s var(--n-bezier);
 box-sizing: content-box;
 display: inline-flex;
 align-items: center;
 justify-content: center;
 `),d(`arrow`,`
 grid-area: arrow;
 font-size: 16px;
 color: var(--n-arrow-color);
 transform: rotate(180deg);
 opacity: 1;
 transition:
 color .3s var(--n-bezier),
 transform 0.2s var(--n-bezier),
 opacity 0.2s var(--n-bezier);
 `),y(`menu-item-content-header`,`
 grid-area: content;
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 opacity: 1;
 white-space: nowrap;
 color: var(--n-item-text-color);
 `,[u(`a`,`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[u(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),d(`extra`,`
 font-size: .93em;
 color: var(--n-group-text-color);
 transition: color .3s var(--n-bezier);
 `)])]),y(`submenu`,`
 cursor: pointer;
 position: relative;
 margin-top: 6px;
 `,[y(`menu-item-content`,`
 height: var(--n-item-height);
 `),y(`submenu-children`,`
 overflow: hidden;
 padding: 0;
 `,[ae({duration:`.2s`})])]),y(`menu-item-group`,[y(`menu-item-group-title`,`
 margin-top: 6px;
 color: var(--n-group-text-color);
 cursor: default;
 font-size: .93em;
 height: 36px;
 display: flex;
 align-items: center;
 transition:
 padding-left .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)])]),y(`menu-tooltip`,[u(`a`,`
 color: inherit;
 text-decoration: none;
 `)]),y(`menu-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function U(e,t){return[_(`hover`,e,t),u(`&:hover`,e,t)]}var W=r({name:`MenuDivider`,setup(){let{mergedClsPrefixRef:e,isHorizontalRef:t}=c(R);return()=>t.value?null:(i(),o(`div`,{key:1,class:p(`${e.value}-menu-divider`)},null,2))}}),G=r({name:`ChevronDownFilled`,render(){return(()=>{let e=m(`f3af82a2aab086a5`);return e[0]||=g(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[g(`path`,{d:`M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z`,fill:`currentColor`})],-1)})()}}),K=[`onClick`],q=r({name:`MenuOptionContent`,props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(e){let{props:t}=c(R);return{menuProps:t,style:s(()=>{let{paddingLeft:t}=e;return{paddingLeft:t&&`${t}px`}}),iconStyle:s(()=>{let{maxIconSize:t,activeIconSize:n,iconMarginRight:r}=e;return{width:`${t}px`,height:`${t}px`,fontSize:`${n}px`,marginRight:`${r}px`}})}},render(){let{clsPrefix:t,tmNode:n,menuProps:{renderIcon:r,renderLabel:a,renderExtra:s,expandIcon:c}}=this,l=r?r(n.rawNode):I(this.icon);return(()=>{let r=m(`7bb10afc6caf8fa4`);return i(),o(`div`,{onClick:e=>{this.onClick?.(e)},role:`none`,class:p([`${t}-menu-item-content`,{[`${t}-menu-item-content--selected`]:this.selected,[`${t}-menu-item-content--collapsed`]:this.collapsed,[`${t}-menu-item-content--child-active`]:this.childActive,[`${t}-menu-item-content--disabled`]:this.disabled,[`${t}-menu-item-content--hover`]:this.hover}]),style:b(this.style)},[T(()=>l&&(i(),o(`div`,{class:p(`${t}-menu-item-content__icon`),style:b(this.iconStyle),role:`none`},[T(()=>[l])],6))),g(`div`,{class:p(`${t}-menu-item-content-header`),role:`none`},[this.isEllipsisPlaceholder?(i(),o(w,{key:0},[T(()=>this.title)],64)):(i(),o(w,{key:1},[a?(i(),o(w,{key:0},[T(()=>a(n.rawNode))],64)):(i(),o(w,{key:1},[T(()=>I(this.title))],64))],64)),this.extra||s?(i(),o(`span`,{key:2,class:p(`${t}-menu-item-content-header__extra`)},[r[0]||=T(` `,-1),s?(i(),o(w,{key:0},[T(()=>s(n.rawNode))],64)):(i(),o(w,{key:1},[T(()=>I(this.extra))],64))],2)):T(()=>null)],2),this.showArrow?(i(),e(h,{key:0,ariaHidden:!0,class:p(`${t}-menu-item-content__arrow`),clsPrefix:t},{default:()=>c?c(n.rawNode):(i(),e(G,{key:1}))},1032,[`class`,`clsPrefix`])):T(()=>null)],14,K)})()}}),le=8;function J(e){let t=c(R),{props:n,mergedCollapsedRef:r}=t,i=c(z,null),a=c(B,null),o=s(()=>n.mode===`horizontal`),l=s(()=>o.value?n.dropdownPlacement:`tmNodes`in e?`right-start`:`right`),u=s(()=>Math.max(n.collapsedIconSize??n.iconSize,n.iconSize));return{dropdownPlacement:l,activeIconSize:s(()=>!o.value&&e.root&&r.value?n.collapsedIconSize??n.iconSize:n.iconSize),maxIconSize:u,paddingLeft:s(()=>{if(o.value)return;let{collapsedWidth:t,indent:s,rootIndent:c}=n,{root:l,isGroup:d}=e,f=c===void 0?s:c;return l?r.value?t/2-u.value/2:f:a&&typeof a.paddingLeftRef.value==`number`?r.value?t/2-u.value/2:s/2+a.paddingLeftRef.value:i&&typeof i.paddingLeftRef.value==`number`?(d?s/2:s)+i.paddingLeftRef.value:0}),iconMarginRight:s(()=>{let{collapsedWidth:t,indent:i,rootIndent:a}=n,{value:s}=u,{root:c}=e;return o.value||!c||!r.value?le:(a===void 0?i:a)+s+le-(t+s)/2}),NMenu:t,NSubmenu:i,NMenuOptionGroup:a}}var Y={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},ue={...Y,tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function},de=O(ue),fe=r({name:`MenuOption`,props:ue,setup(e){let t=J(e),{NSubmenu:n,NMenu:r,NMenuOptionGroup:i}=t,{props:a,mergedClsPrefixRef:o,mergedCollapsedRef:c}=r,l=n?n.mergedDisabledRef:i?i.mergedDisabledRef:{value:!1},u=s(()=>l.value||e.disabled);function d(t){let{onClick:n}=e;n&&n(t)}function f(t){u.value||(r.doSelect(e.internalKey,e.tmNode.rawNode),d(t))}return{mergedClsPrefix:o,dropdownPlacement:t.dropdownPlacement,paddingLeft:t.paddingLeft,iconMarginRight:t.iconMarginRight,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,mergedTheme:r.mergedThemeRef,menuProps:a,dropdownEnabled:M(()=>e.root&&c.value&&a.mode!==`horizontal`&&!u.value),selected:M(()=>r.mergedValueRef.value===e.internalKey),mergedDisabled:u,handleClick:f}},render(){let{mergedClsPrefix:n,mergedTheme:r,tmNode:a,menuProps:{renderLabel:s,nodeProps:c}}=this,l=c?.(a.rawNode);return i(),o(`div`,t(l,{role:`menuitem`,class:[`${n}-menu-item`,l?.class]}),[(i(),e(P,{theme:r.peers.Tooltip,themeOverrides:r.peerOverrides.Tooltip,trigger:`hover`,placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:[`menu-tooltip`]},{default:()=>s?s(a.rawNode):I(this.title),trigger:()=>(i(),e(q,{tmNode:a,clsPrefix:n,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick},null,8,[`tmNode`,`clsPrefix`,`paddingLeft`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`selected`,`title`,`extra`,`disabled`,`icon`,`onClick`]))},1032,[`theme`,`themeOverrides`,`placement`,`disabled`]))],16)}}),X={...Y,tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}},pe=O(X),me=r({name:`MenuOptionGroup`,props:X,setup(e){let n=J(e),{NSubmenu:r}=n,a=s(()=>r?.mergedDisabledRef.value?!0:e.tmNode.disabled);l(B,{paddingLeftRef:n.paddingLeft,mergedDisabledRef:a});let{mergedClsPrefixRef:u,props:d}=c(R);return function(){let{value:r}=u,a=n.paddingLeft.value,{nodeProps:s}=d,c=s?.(e.tmNode.rawNode);return(()=>{let n=m(`45eca6a63be5028b`);return i(),o(`div`,{class:p(`${r}-menu-item-group`),role:`group`},[g(`div`,t(c,{class:[`${r}-menu-item-group-title`,c?.class],style:[c?.style||``,a===void 0?``:`padding-left: ${a}px;`]}),[T(()=>I(e.title)),e.extra?(i(),o(w,{key:0},[n[0]||=T(` `,-1),T(()=>I(e.extra))],64)):T(()=>null)],16),g(`div`,null,[T(()=>e.tmNodes.map(e=>$(e,d)))])],2)})()}}}),he=[`aria-expanded`,`id`],ge=[`aria-expanded`,`id`],_e={...Y,rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean},ve=O(_e),Z=r({name:`Submenu`,props:_e,setup(e){let t=J(e),{NMenu:n,NSubmenu:r}=t,{props:i,mergedCollapsedRef:a,mergedThemeRef:o}=n,c=s(()=>{let{disabled:t}=e;return r?.mergedDisabledRef.value||i.disabled?!0:t}),u=S(!1);l(z,{paddingLeftRef:t.paddingLeft,mergedDisabledRef:c}),l(B,null);function d(){let{onClick:t}=e;t&&t()}function f(){c.value||(a.value||n.toggleExpand(e.internalKey),d())}function p(e){u.value=e}return{menuProps:i,mergedTheme:o,doSelect:n.doSelect,inverted:n.invertedRef,isHorizontal:n.isHorizontalRef,mergedClsPrefix:n.mergedClsPrefixRef,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,iconMarginRight:t.iconMarginRight,dropdownPlacement:t.dropdownPlacement,dropdownShow:u,paddingLeft:t.paddingLeft,mergedDisabled:c,mergedValue:n.mergedValueRef,childActive:M(()=>e.virtualChildActive??n.activePathRef.value.includes(e.internalKey)),collapsed:s(()=>i.mode===`horizontal`?!1:a.value?!0:!n.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:s(()=>!c.value&&(i.mode===`horizontal`||a.value)),handlePopoverShowChange:p,handleClick:f}},render(){let{mergedClsPrefix:n,menuProps:{renderIcon:r,renderLabel:a}}=this,s=()=>{let{isHorizontal:n,paddingLeft:r,collapsed:a,mergedDisabled:s,maxIconSize:c,activeIconSize:l,title:u,childActive:d,icon:f,handleClick:p,menuProps:{nodeProps:m},dropdownShow:h,iconMarginRight:g,tmNode:_,mergedClsPrefix:v,isEllipsisPlaceholder:y,extra:b}=this,x=m?.(_.rawNode);return i(),o(`div`,t(x,{class:[`${v}-menu-item`,x?.class],role:`menuitem`}),[(i(),e(q,{tmNode:_,paddingLeft:r,collapsed:a,disabled:s,iconMarginRight:g,maxIconSize:c,activeIconSize:l,title:u,extra:b,showArrow:!n,childActive:d,clsPrefix:v,icon:f,hover:h,onClick:p,isEllipsisPlaceholder:y},null,8,[`tmNode`,`paddingLeft`,`collapsed`,`disabled`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`title`,`extra`,`showArrow`,`childActive`,`clsPrefix`,`icon`,`hover`,`onClick`,`isEllipsisPlaceholder`]))],16)},c=()=>(i(),e(N,null,{default:()=>{let{tmNodes:e,collapsed:t}=this;return t?null:(i(),o(`div`,{key:1,class:p(`${n}-submenu-children`),role:`menu`},[T(()=>e.map(e=>$(e,this.menuProps)))],2))}},1024));return this.root?(i(),e(F,t({key:2,size:`large`,trigger:`hover`},this.menuProps?.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:`14px`,optionIconSizeLarge:`18px`},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:r,renderLabel:a}),{default:()=>(i(),o(`div`,{class:p(`${n}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[T(()=>s()),this.isHorizontal?T(()=>null):(i(),o(w,{key:1},[T(()=>c())],64))],10,he))},1040,[`themeOverrides`,`theme`,`value`,`disabled`,`placement`,`keyField`,`labelField`,`childrenField`,`onUpdateShow`,`options`,`onSelect`,`inverted`,`renderIcon`,`renderLabel`])):(i(),o(`div`,{key:3,class:p(`${n}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[T(()=>s()),T(()=>c())],10,ge))}});function Q(e){return e.type===`divider`||e.type===`render`}function ye(e){return e.type===`divider`}function $(n,r){let{rawNode:o}=n,{show:s}=o;if(s===!1)return null;if(Q(o))return ye(o)?(i(),e(W,t({key:n.key},o.props),null,16)):null;let{labelField:c}=r,{key:l,level:u,isGroup:d}=n,f={...o,title:o.title||o[c],extra:o.titleExtra||o.extra,key:l,internalKey:l,level:u,root:u===0,isGroup:d};return n.children?n.isGroup?a(me,j(f,pe,{tmNode:n,tmNodes:n.children,key:l})):a(Z,j(f,ve,{key:l,rawNodes:o[r.childrenField],tmNodes:n.children,tmNode:n})):a(fe,j(f,de,{key:l,tmNode:n}))}var be={...v.props,options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},disabledField:{type:String,default:`disabled`},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:`vertical`},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:`bottom`},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array},xe=r({name:`Menu`,inheritAttrs:!1,props:be,setup(t){let{mergedClsPrefixRef:n,inlineThemeDisabled:r}=te(t),a=v(`Menu`,`-menu`,ce,oe,t,n),o=c(se,null),u=s(()=>{let{collapsed:e}=t;if(e!==void 0)return e;if(o){let{collapseModeRef:e,collapsedRef:t}=o;if(e.value===`width`)return t.value??!1}return!1}),d=s(()=>{let{keyField:e,childrenField:n,disabledField:r}=t;return D(t.items||t.options,{getIgnored(e){return Q(e)},getChildren(e){return e[n]},getDisabled(e){return e[r]},getKey(t){return t[e]??t.name}})}),p=s(()=>new Set(d.value.treeNodes.map(e=>e.key))),{watchProps:m}=t,h=S(null);m?.includes(`defaultValue`)?f(()=>{h.value=t.defaultValue}):h.value=t.defaultValue;let g=ee(t,`value`),_=re(g,h),y=S([]),b=()=>{y.value=t.defaultExpandAll?d.value.getNonLeafKeys():t.defaultExpandedNames||t.defaultExpandedKeys||d.value.getPath(_.value,{includeSelf:!1}).keyPath};m?.includes(`defaultExpandedKeys`)?f(b):b();let C=ne(t,[`expandedNames`,`expandedKeys`]),w=re(C,y),T=s(()=>d.value.treeNodes),E=s(()=>d.value.getPath(_.value).keyPath);l(R,{props:t,mergedCollapsedRef:u,mergedThemeRef:a,mergedValueRef:_,mergedExpandedKeysRef:w,activePathRef:E,mergedClsPrefixRef:n,isHorizontalRef:s(()=>t.mode===`horizontal`),invertedRef:ee(t,`inverted`),doSelect:O,toggleExpand:j});function O(e,n){let{"onUpdate:value":r,onUpdateValue:i,onSelect:a}=t;i&&A(i,e,n),r&&A(r,e,n),a&&A(a,e,n),h.value=e}function k(e){let{"onUpdate:expandedKeys":n,onUpdateExpandedKeys:r,onExpandedNamesChange:i,onOpenNamesChange:a}=t;n&&A(n,e),r&&A(r,e),i&&A(i,e),a&&A(a,e),y.value=e}function j(e){let n=Array.from(w.value),r=n.findIndex(t=>t===e);if(~r)n.splice(r,1);else{if(t.accordion&&p.value.has(e)){let e=n.findIndex(e=>p.value.has(e));e>-1&&n.splice(e,1)}n.push(e)}k(n)}let M=e=>{let n=d.value.getPath(e??_.value,{includeSelf:!1}).keyPath;if(!n.length)return;let r=Array.from(w.value),i=new Set([...r,...n]);t.accordion&&p.value.forEach(e=>{i.has(e)&&!n.includes(e)&&i.delete(e)}),k(Array.from(i))},N=s(()=>{let{inverted:e}=t,{common:{cubicBezierEaseInOut:n},self:r}=a.value,{borderRadius:i,borderColorHorizontal:o,fontSize:s,itemHeight:c,dividerColor:l}=r,u={"--n-divider-color":l,"--n-bezier":n,"--n-font-size":s,"--n-border-color-horizontal":o,"--n-border-radius":i,"--n-item-height":c};return e?(u[`--n-group-text-color`]=r.groupTextColorInverted,u[`--n-color`]=r.colorInverted,u[`--n-item-text-color`]=r.itemTextColorInverted,u[`--n-item-text-color-hover`]=r.itemTextColorHoverInverted,u[`--n-item-text-color-active`]=r.itemTextColorActiveInverted,u[`--n-item-text-color-child-active`]=r.itemTextColorChildActiveInverted,u[`--n-item-text-color-child-active-hover`]=r.itemTextColorChildActiveInverted,u[`--n-item-text-color-active-hover`]=r.itemTextColorActiveHoverInverted,u[`--n-item-icon-color`]=r.itemIconColorInverted,u[`--n-item-icon-color-hover`]=r.itemIconColorHoverInverted,u[`--n-item-icon-color-active`]=r.itemIconColorActiveInverted,u[`--n-item-icon-color-active-hover`]=r.itemIconColorActiveHoverInverted,u[`--n-item-icon-color-child-active`]=r.itemIconColorChildActiveInverted,u[`--n-item-icon-color-child-active-hover`]=r.itemIconColorChildActiveHoverInverted,u[`--n-item-icon-color-collapsed`]=r.itemIconColorCollapsedInverted,u[`--n-item-text-color-horizontal`]=r.itemTextColorHorizontalInverted,u[`--n-item-text-color-hover-horizontal`]=r.itemTextColorHoverHorizontalInverted,u[`--n-item-text-color-active-horizontal`]=r.itemTextColorActiveHorizontalInverted,u[`--n-item-text-color-child-active-horizontal`]=r.itemTextColorChildActiveHorizontalInverted,u[`--n-item-text-color-child-active-hover-horizontal`]=r.itemTextColorChildActiveHoverHorizontalInverted,u[`--n-item-text-color-active-hover-horizontal`]=r.itemTextColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-horizontal`]=r.itemIconColorHorizontalInverted,u[`--n-item-icon-color-hover-horizontal`]=r.itemIconColorHoverHorizontalInverted,u[`--n-item-icon-color-active-horizontal`]=r.itemIconColorActiveHorizontalInverted,u[`--n-item-icon-color-active-hover-horizontal`]=r.itemIconColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-child-active-horizontal`]=r.itemIconColorChildActiveHorizontalInverted,u[`--n-item-icon-color-child-active-hover-horizontal`]=r.itemIconColorChildActiveHoverHorizontalInverted,u[`--n-arrow-color`]=r.arrowColorInverted,u[`--n-arrow-color-hover`]=r.arrowColorHoverInverted,u[`--n-arrow-color-active`]=r.arrowColorActiveInverted,u[`--n-arrow-color-active-hover`]=r.arrowColorActiveHoverInverted,u[`--n-arrow-color-child-active`]=r.arrowColorChildActiveInverted,u[`--n-arrow-color-child-active-hover`]=r.arrowColorChildActiveHoverInverted,u[`--n-item-color-hover`]=r.itemColorHoverInverted,u[`--n-item-color-active`]=r.itemColorActiveInverted,u[`--n-item-color-active-hover`]=r.itemColorActiveHoverInverted,u[`--n-item-color-active-collapsed`]=r.itemColorActiveCollapsedInverted):(u[`--n-group-text-color`]=r.groupTextColor,u[`--n-color`]=r.color,u[`--n-item-text-color`]=r.itemTextColor,u[`--n-item-text-color-hover`]=r.itemTextColorHover,u[`--n-item-text-color-active`]=r.itemTextColorActive,u[`--n-item-text-color-child-active`]=r.itemTextColorChildActive,u[`--n-item-text-color-child-active-hover`]=r.itemTextColorChildActiveHover,u[`--n-item-text-color-active-hover`]=r.itemTextColorActiveHover,u[`--n-item-icon-color`]=r.itemIconColor,u[`--n-item-icon-color-hover`]=r.itemIconColorHover,u[`--n-item-icon-color-active`]=r.itemIconColorActive,u[`--n-item-icon-color-active-hover`]=r.itemIconColorActiveHover,u[`--n-item-icon-color-child-active`]=r.itemIconColorChildActive,u[`--n-item-icon-color-child-active-hover`]=r.itemIconColorChildActiveHover,u[`--n-item-icon-color-collapsed`]=r.itemIconColorCollapsed,u[`--n-item-text-color-horizontal`]=r.itemTextColorHorizontal,u[`--n-item-text-color-hover-horizontal`]=r.itemTextColorHoverHorizontal,u[`--n-item-text-color-active-horizontal`]=r.itemTextColorActiveHorizontal,u[`--n-item-text-color-child-active-horizontal`]=r.itemTextColorChildActiveHorizontal,u[`--n-item-text-color-child-active-hover-horizontal`]=r.itemTextColorChildActiveHoverHorizontal,u[`--n-item-text-color-active-hover-horizontal`]=r.itemTextColorActiveHoverHorizontal,u[`--n-item-icon-color-horizontal`]=r.itemIconColorHorizontal,u[`--n-item-icon-color-hover-horizontal`]=r.itemIconColorHoverHorizontal,u[`--n-item-icon-color-active-horizontal`]=r.itemIconColorActiveHorizontal,u[`--n-item-icon-color-active-hover-horizontal`]=r.itemIconColorActiveHoverHorizontal,u[`--n-item-icon-color-child-active-horizontal`]=r.itemIconColorChildActiveHorizontal,u[`--n-item-icon-color-child-active-hover-horizontal`]=r.itemIconColorChildActiveHoverHorizontal,u[`--n-arrow-color`]=r.arrowColor,u[`--n-arrow-color-hover`]=r.arrowColorHover,u[`--n-arrow-color-active`]=r.arrowColorActive,u[`--n-arrow-color-active-hover`]=r.arrowColorActiveHover,u[`--n-arrow-color-child-active`]=r.arrowColorChildActive,u[`--n-arrow-color-child-active-hover`]=r.arrowColorChildActiveHover,u[`--n-item-color-hover`]=r.itemColorHover,u[`--n-item-color-active`]=r.itemColorActive,u[`--n-item-color-active-hover`]=r.itemColorActiveHover,u[`--n-item-color-active-collapsed`]=r.itemColorActiveCollapsed),u}),P=r?x(`menu`,s(()=>t.inverted?`a`:`b`),N,t):void 0,F=ie(),I=S(null),ae=S(null),L=!0,z=()=>{L?L=!1:I.value?.sync({showAllItemsBeforeCalculate:!0})};function B(){return document.getElementById(F)}let V=S(-1);function H(e){V.value=t.options.length-e}function U(e){e||(V.value=-1)}let W=s(()=>{let e=V.value;return{children:e===-1?[]:t.options.slice(e)}}),G=s(()=>{let{childrenField:e,disabledField:n,keyField:r}=t;return D([W.value],{getIgnored(e){return Q(e)},getChildren(t){return t[e]},getDisabled(e){return e[n]},getKey(e){return e[r]??e.name}})}),K=s(()=>D([{}]).treeNodes[0]);function q(){if(V.value===-1)return i(),e(Z,{root:!0,level:0,key:`__ellpisisGroupPlaceholder__`,internalKey:`__ellpisisGroupPlaceholder__`,title:`···`,tmNode:K.value,domId:F,isEllipsisPlaceholder:!0},null,8,[`tmNode`,`domId`]);let t=G.value.treeNodes[0],n=E.value,r=!!t.children?.some(e=>n.includes(e.key));return i(),e(Z,{level:0,root:!0,key:`__ellpisisGroup__`,internalKey:`__ellpisisGroup__`,title:`···`,virtualChildActive:r,tmNode:t,domId:F,rawNodes:t.rawNode.children||[],tmNodes:t.children||[],isEllipsisPlaceholder:!0},null,8,[`virtualChildActive`,`tmNode`,`domId`,`rawNodes`,`tmNodes`])}return{mergedClsPrefix:n,controlledExpandedKeys:C,uncontrolledExpanededKeys:y,mergedExpandedKeys:w,uncontrolledValue:h,mergedValue:_,activePath:E,tmNodes:T,mergedTheme:a,mergedCollapsed:u,cssVars:r?void 0:N,themeClass:P?.themeClass,overflowRef:I,counterRef:ae,updateCounter:()=>{},onResize:z,onUpdateOverflow:U,onUpdateCount:H,renderCounter:q,getCounter:B,onRender:P?.onRender,showOption:M,deriveResponsiveState:z}},render(){let{mergedClsPrefix:n,mode:r,themeClass:o,onRender:s}=this;s?.();let c=()=>this.tmNodes.map(e=>$(e,this.$props)),l=r===`horizontal`&&this.responsive,u=()=>a(`div`,t(this.$attrs,{role:r===`horizontal`?`menubar`:`menu`,class:[`${n}-menu`,o,`${n}-menu--${r}`,l&&`${n}-menu--responsive`,this.mergedCollapsed&&`${n}-menu--collapsed`],style:this.cssVars}),l?(i(),e(E,{key:2,ref:`overflowRef`,onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:c,counter:this.renderCounter},1032,[`onUpdateOverflow`,`getCounter`,`onUpdateCount`,`updateCounter`])):c());return l?(i(),e(k,{key:3,onResize:this.onResize},{default:u},1032,[`onResize`])):u()}});export{se as n,L as r,xe as t};