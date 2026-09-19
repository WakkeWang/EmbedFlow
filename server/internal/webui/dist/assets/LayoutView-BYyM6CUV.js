import{At as e,Bt as t,D as n,Ft as r,Jt as i,Kt as a,Mt as o,Nt as s,Ot as c,Pt as l,Qt as u,Rt as d,Yt as f,a as p,at as m,c as h,ct as g,dn as _,f as v,g as y,i as b,kt as x,lt as S,mn as C,n as w,nn as ee,o as T,ot as E,pn as D,s as O,sn as k,tn as A,u as j,un as M,v as N,w as P,wt as te,xt as F,y as I,yt as ne}from"./vue-i18n-CgKceEWx.js";import{t as L}from"./format-length-04qCp4Ww.js";import{c as R,m as z,n as B,t as V}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{f as re}from"./fade-in-scale-up.cssr-7akLf3Ad.js";import{i as ie}from"./Dropdown-BlbHn3LV.js";import{t as ae}from"./Select-D7J9B1aO.js";import{l as H,t as U}from"./http-BdILZoCd.js";import{n as W,r as G,t as K}from"./Menu-BbN4CukZ.js";import{i as q,n as J,r as oe}from"./index-kjlF3dpY.js";import{t as Y}from"./project-CUm4GFJ_.js";function se(e){let{baseColor:t,textColor2:n,bodyColor:r,cardColor:i,dividerColor:a,actionColor:o,scrollbarColor:s,scrollbarColorHover:c,invertedColor:l}=e;return{textColor:n,textColorInverted:`#FFF`,color:r,colorEmbedded:o,headerColor:i,headerColorInverted:l,footerColor:o,footerColorInverted:l,headerBorderColor:a,headerBorderColorInverted:l,footerBorderColor:a,footerBorderColorInverted:l,siderBorderColor:a,siderBorderColorInverted:l,siderColor:i,siderColorInverted:l,siderToggleButtonBorder:`1px solid ${a}`,siderToggleButtonColor:t,siderToggleButtonIconColor:n,siderToggleButtonIconColorInverted:n,siderToggleBarColor:v(r,s),siderToggleBarColorHover:v(r,c),__invertScrollbar:`true`}}var X=p({name:`Layout`,common:j,peers:{Scrollbar:h},self:se}),ce=E(`layout`,`
 color: var(--n-text-color);
 background-color: var(--n-color);
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 flex: auto;
 overflow: hidden;
 transition:
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
`,[E(`layout-scroll-container`,`
 overflow-x: hidden;
 box-sizing: border-box;
 height: 100%;
 `),S(`absolute-positioned`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),le={embedded:Boolean,position:G,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:``},hasSider:Boolean,siderPlacement:{type:String,default:`left`}},Z=n(`n-layout`);function Q(n){return r({name:n?`LayoutContent`:`Layout`,props:{...T.props,...le},setup(e){let t=k(null),n=k(null),{mergedClsPrefixRef:r,inlineThemeDisabled:i}=P(e),a=T(`Layout`,`-layout`,ce,X,e,r);function o(r,i){if(e.nativeScrollbar){let{value:e}=t;e&&(i===void 0?e.scrollTo(r):e.scrollTo(r,i))}else{let{value:e}=n;e&&e.scrollTo(r,i)}}f(Z,e);let s=0,l=0,u=t=>{let n=t.target;s=n.scrollLeft,l=n.scrollTop,e.onScroll?.(t)};R(()=>{if(e.nativeScrollbar){let e=t.value;e&&(e.scrollTop=l,e.scrollLeft=s)}});let d={display:`flex`,flexWrap:`nowrap`,width:`100%`,flexDirection:`row`},p={scrollTo:o},m=c(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=a.value;return{"--n-bezier":t,"--n-color":e.embedded?n.colorEmbedded:n.color,"--n-text-color":n.textColor}}),h=i?O(`layout`,c(()=>e.embedded?`e`:``),m,e):void 0;return{mergedClsPrefix:r,scrollableElRef:t,scrollbarInstRef:n,hasSiderStyle:d,mergedTheme:a,handleNativeElScroll:u,cssVars:i?void 0:m,themeClass:h?.themeClass,onRender:h?.onRender,...p}},render(){let{mergedClsPrefix:r,hasSider:a}=this;this.onRender?.();let s=a?this.hasSiderStyle:void 0,c=[this.themeClass,n&&`${r}-layout-content`,`${r}-layout`,`${r}-layout--${this.position}-positioned`];return i(),o(`div`,{class:y(c),style:D(this.cssVars)},[this.nativeScrollbar?(i(),o(`div`,{key:0,ref:`scrollableElRef`,class:y([`${r}-layout-scroll-container`,this.contentClass]),style:D([this.contentStyle,s]),onScroll:this.handleNativeElScroll},[I(()=>this.$slots.default?.())],46,[`onScroll`])):(i(),e(B,t({key:1},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,s]}),N(this.$slots),1040,[`onScroll`,`theme`,`themeOverrides`,`contentClass`,`contentStyle`]))],6)}})}var ue=Q(!1),de=Q(!0),fe=E(`layout-sider`,`
 flex-shrink: 0;
 box-sizing: border-box;
 position: relative;
 z-index: 1;
 color: var(--n-text-color);
 transition:
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 min-width .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 transform .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 display: flex;
 justify-content: flex-end;
`,[S(`bordered`,[g(`border`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),g(`left-placement`,[S(`bordered`,[g(`border`,`
 right: 0;
 `)])]),S(`right-placement`,`
 justify-content: flex-start;
 `,[S(`bordered`,[g(`border`,`
 left: 0;
 `)]),S(`collapsed`,[E(`layout-toggle-button`,[E(`base-icon`,`
 transform: rotate(180deg);
 `)]),E(`layout-toggle-bar`,[m(`&:hover`,[g(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),g(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])])]),E(`layout-toggle-button`,`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[E(`base-icon`,`
 transform: rotate(0);
 `)]),E(`layout-toggle-bar`,`
 left: -28px;
 transform: rotate(180deg);
 `,[m(`&:hover`,[g(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),g(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})])])]),S(`collapsed`,[E(`layout-toggle-bar`,[m(`&:hover`,[g(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),g(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])]),E(`layout-toggle-button`,[E(`base-icon`,`
 transform: rotate(0);
 `)])]),E(`layout-toggle-button`,`
 transition:
 color .3s var(--n-bezier),
 right .3s var(--n-bezier),
 left .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 cursor: pointer;
 width: 24px;
 height: 24px;
 position: absolute;
 top: 50%;
 right: 0;
 border-radius: 50%;
 display: flex;
 align-items: center;
 justify-content: center;
 font-size: 18px;
 color: var(--n-toggle-button-icon-color);
 border: var(--n-toggle-button-border);
 background-color: var(--n-toggle-button-color);
 box-shadow: 0 2px 4px 0px rgba(0, 0, 0, .06);
 transform: translateX(50%) translateY(-50%);
 z-index: 1;
 `,[E(`base-icon`,`
 transition: transform .3s var(--n-bezier);
 transform: rotate(180deg);
 `)]),E(`layout-toggle-bar`,`
 cursor: pointer;
 height: 72px;
 width: 32px;
 position: absolute;
 top: calc(50% - 36px);
 right: -28px;
 `,[g(`top, bottom`,`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),g(`bottom`,`
 position: absolute;
 top: 34px;
 `),m(`&:hover`,[g(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),g(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})]),g(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color)`}),m(`&:hover`,[g(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color-hover)`})])]),g(`border`,`
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 width: 1px;
 transition: background-color .3s var(--n-bezier);
 `),E(`layout-sider-scroll-container`,`
 flex-grow: 1;
 flex-shrink: 0;
 box-sizing: border-box;
 height: 100%;
 opacity: 0;
 transition: opacity .3s var(--n-bezier);
 max-width: 100%;
 `),S(`show-content`,[E(`layout-sider-scroll-container`,{opacity:1})]),S(`absolute-positioned`,`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]),pe=[`onClick`],me=r({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:e}=this;return i(),o(`div`,{onClick:this.onClick,class:y(`${e}-layout-toggle-bar`)},[x(`div`,{class:y(`${e}-layout-toggle-bar__top`)},null,2),x(`div`,{class:y(`${e}-layout-toggle-bar__bottom`)},null,2)],10,pe)}}),$=[`onClick`],he=r({name:`LayoutToggleButton`,props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:t}=this;return i(),o(`div`,{class:y(`${t}-layout-toggle-button`),onClick:this.onClick},[(i(),e(b,{clsPrefix:t},{default:()=>(i(),e(ie))},1032,[`clsPrefix`]))],10,$)}}),ge=[`onTransitionend`],_e={position:G,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:``},collapseMode:{type:String,default:`transform`},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function},ve=r({name:`LayoutSider`,props:{...T.props,..._e},setup(e){let t=d(Z),n=k(null),r=k(null),i=k(e.defaultCollapsed),a=re(M(e,`collapsed`),i),o=c(()=>L(a.value?e.collapsedWidth:e.width)),s=c(()=>e.collapseMode===`transform`?{minWidth:L(e.width)}:{}),l=c(()=>t?t.siderPlacement:`left`);function u(t,i){if(e.nativeScrollbar){let{value:e}=n;e&&(i===void 0?e.scrollTo(t):e.scrollTo(t,i))}else{let{value:e}=r;e&&e.scrollTo(t,i)}}function p(){let{"onUpdate:collapsed":t,onUpdateCollapsed:n,onExpand:r,onCollapse:o}=e,{value:s}=a;n&&z(n,!s),t&&z(t,!s),i.value=!s,s?r&&z(r):o&&z(o)}let m=0,h=0,g=t=>{let n=t.target;m=n.scrollLeft,h=n.scrollTop,e.onScroll?.(t)};R(()=>{if(e.nativeScrollbar){let e=n.value;e&&(e.scrollTop=h,e.scrollLeft=m)}}),f(W,{collapsedRef:a,collapseModeRef:M(e,`collapseMode`)});let{mergedClsPrefixRef:_,inlineThemeDisabled:v}=P(e),y=T(`Layout`,`-layout-sider`,fe,X,e,_);function b(t){t.propertyName===`max-width`&&(a.value?e.onAfterLeave?.():e.onAfterEnter?.())}let x={scrollTo:u},S=c(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=y.value,{siderToggleButtonColor:r,siderToggleButtonBorder:i,siderToggleBarColor:a,siderToggleBarColorHover:o}=n,s={"--n-bezier":t,"--n-toggle-button-color":r,"--n-toggle-button-border":i,"--n-toggle-bar-color":a,"--n-toggle-bar-color-hover":o};return e.inverted?(s[`--n-color`]=n.siderColorInverted,s[`--n-text-color`]=n.textColorInverted,s[`--n-border-color`]=n.siderBorderColorInverted,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColorInverted,s.__invertScrollbar=n.__invertScrollbar):(s[`--n-color`]=n.siderColor,s[`--n-text-color`]=n.textColor,s[`--n-border-color`]=n.siderBorderColor,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColor),s}),C=v?O(`layout-sider`,c(()=>e.inverted?`a`:`b`),S,e):void 0;return{scrollableElRef:n,scrollbarInstRef:r,mergedClsPrefix:_,mergedTheme:y,styleMaxWidth:o,mergedCollapsed:a,scrollContainerStyle:s,siderPlacement:l,handleNativeElScroll:g,handleTransitionend:b,handleTriggerClick:p,inlineThemeDisabled:v,cssVars:S,themeClass:C?.themeClass,onRender:C?.onRender,...x}},render(){let{mergedClsPrefix:n,mergedCollapsed:r,showTrigger:a}=this;return this.onRender?.(),i(),o(`aside`,{class:y([`${n}-layout-sider`,this.themeClass,`${n}-layout-sider--${this.position}-positioned`,`${n}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${n}-layout-sider--bordered`,r&&`${n}-layout-sider--collapsed`,(!r||this.showCollapsedContent)&&`${n}-layout-sider--show-content`]),onTransitionend:this.handleTransitionend,style:D([this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:L(this.width)}])},[this.nativeScrollbar?(i(),o(`div`,{key:1,class:y([`${n}-layout-sider-scroll-container`,this.contentClass]),onScroll:this.handleNativeElScroll,style:D([this.scrollContainerStyle,{overflow:`auto`},this.contentStyle]),ref:`scrollableElRef`},[I(()=>this.$slots.default?.())],46,[`onScroll`])):(i(),e(B,t({key:0},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar===`true`?{colorHover:`rgba(255, 255, 255, .4)`,color:`rgba(255, 255, 255, .3)`}:void 0}),N(this.$slots),1040,[`onScroll`,`style`,`contentStyle`,`contentClass`,`theme`,`themeOverrides`,`builtinThemeOverrides`])),a?(i(),o(te,{key:2},[a===`bar`?(i(),e(me,{key:0,clsPrefix:n,class:y(r?this.collapsedTriggerClass:this.triggerClass),style:D(r?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`])):(i(),e(he,{key:1,clsPrefix:n,class:y(r?this.collapsedTriggerClass:this.triggerClass),style:D(r?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`]))],64)):I(()=>null),this.bordered?(i(),o(`div`,{key:4,class:y(`${n}-layout-sider__border`)},null,2)):I(()=>null)],46,ge)}}),ye={class:`project-switcher`},be={key:1,class:`none-hint`},xe={key:3,class:`quick-create`},Se=[`placeholder`],Ce=V(r({__name:`ProjectSwitcher`,setup(t){let{t:n}=w(),{projects:r,currentId:a,select:c,create:u}=Y(),d=k(!1),f=k(``);async function p(){f.value.trim()&&(await u(f.value.trim(),``),f.value=``,d.value=!1)}return(t,u)=>(i(),o(`div`,ye,[_(r).length>0?(i(),e(_(ae),{key:0,value:_(a),options:_(r).map(e=>({label:e.name,value:e.id})),size:`small`,placeholder:_(n)(`nav.projectNone`),"onUpdate:value":u[0]||=e=>_(c)(e)},null,8,[`value`,`options`,`placeholder`])):(i(),o(`span`,be,C(_(n)(`nav.projectNone`)),1)),d.value?(i(),o(`div`,xe,[ee(x(`input`,{"onUpdate:modelValue":u[2]||=e=>f.value=e,class:`n-input`,placeholder:_(n)(`project.name`),onKeyup:F(p,[`enter`])},null,40,Se),[[ne,f.value]]),l(_(H),{size:`tiny`,type:`primary`,onClick:p},{default:A(()=>[s(C(_(n)(`project.create`)),1)]),_:1})])):(i(),e(_(H),{key:2,quaternary:``,size:`tiny`,onClick:u[1]||=e=>d.value=!0},{default:A(()=>[...u[3]||=[s(`+`,-1)]]),_:1}))]))}}),[[`__scopeId`,`data-v-ab0600dc`]]),we={class:`sider-footer`},Te=V(r({__name:`LayoutView`,setup(t){let{t:n,locale:r}=w(),o=J(),{load:d}=Y();a(d);let f=c(()=>[{label:n(`nav.dashboard`),key:`/`},{label:n(`nav.build`),key:`/build`,disabled:!1},{label:n(`nav.deploy`),key:`/deploy`},{label:n(`nav.test`),key:`/test`,disabled:!1},{label:n(`nav.release`),key:`/release`,disabled:!1},{label:n(`nav.settings`),key:`/settings`,disabled:!0}]);function p(e){o.push(e)}function m(){r.value=r.value===`zh`?`en`:`zh`}function h(){U(),o.push(`/login`)}return(t,r)=>{let a=u(`router-view`);return i(),e(_(ue),{position:`absolute`,class:`root`},{default:A(()=>[l(_(ve),{bordered:``,width:200,class:`sider`},{default:A(()=>[r[0]||=x(`div`,{class:`brand`},`EmbedFlow`,-1),l(Ce),l(_(K),{options:f.value,"onUpdate:value":p},null,8,[`options`]),x(`div`,we,[l(_(H),{quaternary:``,size:`small`,onClick:_(q)},{default:A(()=>[s(C(_(oe).isDark?`Light`:`Dark`),1)]),_:1},8,[`onClick`]),l(_(H),{quaternary:``,size:`small`,onClick:m},{default:A(()=>[s(C(_(n)(`nav.language`)),1)]),_:1}),l(_(H),{quaternary:``,size:`small`,onClick:h},{default:A(()=>[s(C(_(n)(`nav.logout`)),1)]),_:1})])]),_:1}),l(_(de),{"content-style":`padding: 16px;`,class:`content`},{default:A(()=>[l(a)]),_:1})]),_:1})}}}),[[`__scopeId`,`data-v-f0c9a85c`]]);export{Te as default};