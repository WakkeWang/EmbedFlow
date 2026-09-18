import{Bt as e,E as t,Gt as n,Ht as r,It as i,Jt as a,Kt as o,L as s,Lt as c,Mt as l,O as u,P as d,Rt as f,S as p,Sn as m,T as h,Ut as g,Vt as _,bt as v,c as y,cn as b,d as x,f as S,g as C,gt as w,h as T,hn as E,ht as D,k as O,ln as k,m as A,n as j,nn as M,on as N,p as ee,rn as P,v as te,vn as F,vt as I,xn as L,y as ne,yn as R,yt as z}from"./vue-i18n-BxOqVJSv.js";import{o as re,p as ie,t as B}from"./create-CKNSf1FL.js";import{t as V}from"./format-length-CqnQVDgE.js";import{I as H,O as U,S as W,g as G,l as K,n as ae,t as oe,v as se}from"./_plugin-vue_export-helper-3sd4LcYN.js";import{c as q,f as J}from"./fade-in-scale-up.cssr-DpVYUTqJ.js";import{i as ce,n as le,t as ue}from"./Dropdown-DtemoHYo.js";import{E as Y,O as de,i as fe,n as pe,o as me,r as he}from"./index-0G3VD5Ph.js";function ge(e){let{baseColor:t,textColor2:n,bodyColor:r,cardColor:i,dividerColor:a,actionColor:o,scrollbarColor:s,scrollbarColorHover:c,invertedColor:l}=e;return{textColor:n,textColorInverted:`#FFF`,color:r,colorEmbedded:o,headerColor:i,headerColorInverted:l,footerColor:o,footerColorInverted:l,headerBorderColor:a,headerBorderColorInverted:l,footerBorderColor:a,footerBorderColorInverted:l,siderBorderColor:a,siderBorderColorInverted:l,siderColor:i,siderColorInverted:l,siderToggleButtonBorder:`1px solid ${a}`,siderToggleButtonColor:t,siderToggleButtonIconColor:n,siderToggleButtonIconColorInverted:n,siderToggleBarColor:p(r,s),siderToggleBarColorHover:p(r,c),__invertScrollbar:`true`}}var _e=ee({name:`Layout`,common:te,peers:{Scrollbar:C},self:ge}),ve=s(`n-layout-sider`),ye={type:String,default:`static`},be=w(`layout`,`
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
`,[w(`layout-scroll-container`,`
 overflow-x: hidden;
 box-sizing: border-box;
 height: 100%;
 `),z(`absolute-positioned`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),xe={embedded:Boolean,position:ye,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:``},hasSider:Boolean,siderPlacement:{type:String,default:`left`}},Se=s(`n-layout`);function Ce(n){return g({name:n?`LayoutContent`:`Layout`,props:{...A.props,...xe},setup(e){let t=E(null),n=E(null),{mergedClsPrefixRef:r,inlineThemeDisabled:a}=d(e),o=A(`Layout`,`-layout`,be,_e,e,r);function s(r,i){if(e.nativeScrollbar){let{value:e}=t;e&&(i===void 0?e.scrollTo(r):e.scrollTo(r,i))}else{let{value:e}=n;e&&e.scrollTo(r,i)}}P(Se,e);let c=0,l=0,u=t=>{let n=t.target;c=n.scrollLeft,l=n.scrollTop,e.onScroll?.(t)};W(()=>{if(e.nativeScrollbar){let e=t.value;e&&(e.scrollTop=l,e.scrollLeft=c)}});let f={display:`flex`,flexWrap:`nowrap`,width:`100%`,flexDirection:`row`},p={scrollTo:s},m=i(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=o.value;return{"--n-bezier":t,"--n-color":e.embedded?n.colorEmbedded:n.color,"--n-text-color":n.textColor}}),h=a?T(`layout`,i(()=>e.embedded?`e`:``),m,e):void 0;return{mergedClsPrefix:r,scrollableElRef:t,scrollbarInstRef:n,hasSiderStyle:f,mergedTheme:o,handleNativeElScroll:u,cssVars:a?void 0:m,themeClass:h?.themeClass,onRender:h?.onRender,...p}},render(){let{mergedClsPrefix:r,hasSider:i}=this;this.onRender?.();let o=i?this.hasSiderStyle:void 0,s=[this.themeClass,n&&`${r}-layout-content`,`${r}-layout`,`${r}-layout--${this.position}-positioned`];return M(),e(`div`,{class:t(s),style:L(this.cssVars)},[this.nativeScrollbar?(M(),e(`div`,{key:0,ref:`scrollableElRef`,class:t([`${r}-layout-scroll-container`,this.contentClass]),style:L([this.contentStyle,o]),onScroll:this.handleNativeElScroll},[O(()=>this.$slots.default?.())],46,[`onScroll`])):(M(),f(G,a({key:1},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,o]}),u(this.$slots),1040,[`onScroll`,`theme`,`themeOverrides`,`contentClass`,`contentStyle`]))],6)}})}var we=Ce(!1),Te=Ce(!0),Ee=w(`layout-sider`,`
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
`,[z(`bordered`,[I(`border`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),I(`left-placement`,[z(`bordered`,[I(`border`,`
 right: 0;
 `)])]),z(`right-placement`,`
 justify-content: flex-start;
 `,[z(`bordered`,[I(`border`,`
 left: 0;
 `)]),z(`collapsed`,[w(`layout-toggle-button`,[w(`base-icon`,`
 transform: rotate(180deg);
 `)]),w(`layout-toggle-bar`,[D(`&:hover`,[I(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),I(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])])]),w(`layout-toggle-button`,`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[w(`base-icon`,`
 transform: rotate(0);
 `)]),w(`layout-toggle-bar`,`
 left: -28px;
 transform: rotate(180deg);
 `,[D(`&:hover`,[I(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),I(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})])])]),z(`collapsed`,[w(`layout-toggle-bar`,[D(`&:hover`,[I(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),I(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])]),w(`layout-toggle-button`,[w(`base-icon`,`
 transform: rotate(0);
 `)])]),w(`layout-toggle-button`,`
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
 `,[w(`base-icon`,`
 transition: transform .3s var(--n-bezier);
 transform: rotate(180deg);
 `)]),w(`layout-toggle-bar`,`
 cursor: pointer;
 height: 72px;
 width: 32px;
 position: absolute;
 top: calc(50% - 36px);
 right: -28px;
 `,[I(`top, bottom`,`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),I(`bottom`,`
 position: absolute;
 top: 34px;
 `),D(`&:hover`,[I(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),I(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})]),I(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color)`}),D(`&:hover`,[I(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color-hover)`})])]),I(`border`,`
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 width: 1px;
 transition: background-color .3s var(--n-bezier);
 `),w(`layout-sider-scroll-container`,`
 flex-grow: 1;
 flex-shrink: 0;
 box-sizing: border-box;
 height: 100%;
 opacity: 0;
 transition: opacity .3s var(--n-bezier);
 max-width: 100%;
 `),z(`show-content`,[w(`layout-sider-scroll-container`,{opacity:1})]),z(`absolute-positioned`,`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]),De=[`onClick`],Oe=g({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:n}=this;return M(),e(`div`,{onClick:this.onClick,class:t(`${n}-layout-toggle-bar`)},[c(`div`,{class:t(`${n}-layout-toggle-bar__top`)},null,2),c(`div`,{class:t(`${n}-layout-toggle-bar__bottom`)},null,2)],10,De)}}),ke=[`onClick`],Ae=g({name:`LayoutToggleButton`,props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:n}=this;return M(),e(`div`,{class:t(`${n}-layout-toggle-button`),onClick:this.onClick},[(M(),f(S,{clsPrefix:n},{default:()=>(M(),f(ce))},1032,[`clsPrefix`]))],10,ke)}}),je=[`onTransitionend`],Me={position:ye,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:``},collapseMode:{type:String,default:`transform`},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function},Ne=g({name:`LayoutSider`,props:{...A.props,...Me},setup(e){let t=o(Se),n=E(null),r=E(null),a=E(e.defaultCollapsed),s=J(F(e,`collapsed`),a),c=i(()=>V(s.value?e.collapsedWidth:e.width)),l=i(()=>e.collapseMode===`transform`?{minWidth:V(e.width)}:{}),u=i(()=>t?t.siderPlacement:`left`);function f(t,i){if(e.nativeScrollbar){let{value:e}=n;e&&(i===void 0?e.scrollTo(t):e.scrollTo(t,i))}else{let{value:e}=r;e&&e.scrollTo(t,i)}}function p(){let{"onUpdate:collapsed":t,onUpdateCollapsed:n,onExpand:r,onCollapse:i}=e,{value:o}=s;n&&U(n,!o),t&&U(t,!o),a.value=!o,o?r&&U(r):i&&U(i)}let m=0,h=0,g=t=>{let n=t.target;m=n.scrollLeft,h=n.scrollTop,e.onScroll?.(t)};W(()=>{if(e.nativeScrollbar){let e=n.value;e&&(e.scrollTop=h,e.scrollLeft=m)}}),P(ve,{collapsedRef:s,collapseModeRef:F(e,`collapseMode`)});let{mergedClsPrefixRef:_,inlineThemeDisabled:v}=d(e),y=A(`Layout`,`-layout-sider`,Ee,_e,e,_);function b(t){t.propertyName===`max-width`&&(s.value?e.onAfterLeave?.():e.onAfterEnter?.())}let x={scrollTo:f},S=i(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=y.value,{siderToggleButtonColor:r,siderToggleButtonBorder:i,siderToggleBarColor:a,siderToggleBarColorHover:o}=n,s={"--n-bezier":t,"--n-toggle-button-color":r,"--n-toggle-button-border":i,"--n-toggle-bar-color":a,"--n-toggle-bar-color-hover":o};return e.inverted?(s[`--n-color`]=n.siderColorInverted,s[`--n-text-color`]=n.textColorInverted,s[`--n-border-color`]=n.siderBorderColorInverted,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColorInverted,s.__invertScrollbar=n.__invertScrollbar):(s[`--n-color`]=n.siderColor,s[`--n-text-color`]=n.textColor,s[`--n-border-color`]=n.siderBorderColor,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColor),s}),C=v?T(`layout-sider`,i(()=>e.inverted?`a`:`b`),S,e):void 0;return{scrollableElRef:n,scrollbarInstRef:r,mergedClsPrefix:_,mergedTheme:y,styleMaxWidth:c,mergedCollapsed:s,scrollContainerStyle:l,siderPlacement:u,handleNativeElScroll:g,handleTransitionend:b,handleTriggerClick:p,inlineThemeDisabled:v,cssVars:S,themeClass:C?.themeClass,onRender:C?.onRender,...x}},render(){let{mergedClsPrefix:n,mergedCollapsed:r,showTrigger:i}=this;return this.onRender?.(),M(),e(`aside`,{class:t([`${n}-layout-sider`,this.themeClass,`${n}-layout-sider--${this.position}-positioned`,`${n}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${n}-layout-sider--bordered`,r&&`${n}-layout-sider--collapsed`,(!r||this.showCollapsedContent)&&`${n}-layout-sider--show-content`]),onTransitionend:this.handleTransitionend,style:L([this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:V(this.width)}])},[this.nativeScrollbar?(M(),e(`div`,{key:1,class:t([`${n}-layout-sider-scroll-container`,this.contentClass]),onScroll:this.handleNativeElScroll,style:L([this.scrollContainerStyle,{overflow:`auto`},this.contentStyle]),ref:`scrollableElRef`},[O(()=>this.$slots.default?.())],46,[`onScroll`])):(M(),f(G,a({key:0},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar===`true`?{colorHover:`rgba(255, 255, 255, .4)`,color:`rgba(255, 255, 255, .3)`}:void 0}),u(this.$slots),1040,[`onScroll`,`style`,`contentStyle`,`contentClass`,`theme`,`themeOverrides`,`builtinThemeOverrides`])),i?(M(),e(l,{key:2},[i===`bar`?(M(),f(Oe,{key:0,clsPrefix:n,class:t(r?this.collapsedTriggerClass:this.triggerClass),style:L(r?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`])):(M(),f(Ae,{key:1,clsPrefix:n,class:t(r?this.collapsedTriggerClass:this.triggerClass),style:L(r?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`]))],64)):O(()=>null),this.bordered?(M(),e(`div`,{key:4,class:t(`${n}-layout-sider__border`)},null,2)):O(()=>null)],46,je)}}),X=s(`n-menu`),Pe=s(`n-submenu`),Z=s(`n-menu-item-group`),Fe=[D(`&::before`,`background-color: var(--n-item-color-hover);`),I(`arrow`,`
 color: var(--n-arrow-color-hover);
 `),I(`icon`,`
 color: var(--n-item-icon-color-hover);
 `),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover);
 `,[D(`a`,`
 color: var(--n-item-text-color-hover);
 `),I(`extra`,`
 color: var(--n-item-text-color-hover);
 `)])],Ie=[I(`icon`,`
 color: var(--n-item-icon-color-hover-horizontal);
 `),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover-horizontal);
 `,[D(`a`,`
 color: var(--n-item-text-color-hover-horizontal);
 `),I(`extra`,`
 color: var(--n-item-text-color-hover-horizontal);
 `)])],Le=D([w(`menu`,`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[z(`horizontal`,`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[w(`submenu`,`margin: 0;`),w(`menu-item`,`margin: 0;`),w(`menu-item-content`,`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[D(`&::before`,`display: none;`),z(`selected`,`border-bottom: 2px solid var(--n-border-color-horizontal)`)]),w(`menu-item-content`,[z(`selected`,[I(`icon`,`color: var(--n-item-icon-color-active-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-horizontal);
 `,[D(`a`,`color: var(--n-item-text-color-active-horizontal);`),I(`extra`,`color: var(--n-item-text-color-active-horizontal);`)])]),z(`child-active`,`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[D(`a`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `),I(`extra`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),I(`icon`,`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),v(`disabled`,[v(`selected, child-active`,[D(`&:focus-within`,Ie)]),z(`selected`,[Q(null,[I(`icon`,`color: var(--n-item-icon-color-active-hover-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[D(`a`,`color: var(--n-item-text-color-active-hover-horizontal);`),I(`extra`,`color: var(--n-item-text-color-active-hover-horizontal);`)])])]),z(`child-active`,[Q(null,[I(`icon`,`color: var(--n-item-icon-color-child-active-hover-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[D(`a`,`color: var(--n-item-text-color-child-active-hover-horizontal);`),I(`extra`,`color: var(--n-item-text-color-child-active-hover-horizontal);`)])])]),Q(`border-bottom: 2px solid var(--n-border-color-horizontal);`,Ie)]),w(`menu-item-content-header`,[D(`a`,`color: var(--n-item-text-color-horizontal);`)])])]),v(`responsive`,[w(`menu-item-content-header`,`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),z(`collapsed`,[w(`menu-item-content`,[z(`selected`,[D(`&::before`,`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),w(`menu-item-content-header`,`opacity: 0;`),I(`arrow`,`opacity: 0;`),I(`icon`,`color: var(--n-item-icon-color-collapsed);`)])]),w(`menu-item`,`
 height: var(--n-item-height);
 margin-top: 6px;
 position: relative;
 `),w(`menu-item-content`,`
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
 `,[D(`> *`,`z-index: 1;`),D(`&::before`,`
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
 `),z(`disabled`,`
 opacity: .45;
 cursor: not-allowed;
 `),z(`collapsed`,[I(`arrow`,`transform: rotate(0);`)]),z(`selected`,[D(`&::before`,`background-color: var(--n-item-color-active);`),I(`arrow`,`color: var(--n-arrow-color-active);`),I(`icon`,`color: var(--n-item-icon-color-active);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active);
 `,[D(`a`,`color: var(--n-item-text-color-active);`),I(`extra`,`color: var(--n-item-text-color-active);`)])]),z(`child-active`,[w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active);
 `,[D(`a`,`
 color: var(--n-item-text-color-child-active);
 `),I(`extra`,`
 color: var(--n-item-text-color-child-active);
 `)]),I(`arrow`,`
 color: var(--n-arrow-color-child-active);
 `),I(`icon`,`
 color: var(--n-item-icon-color-child-active);
 `)]),v(`disabled`,[v(`selected, child-active`,[D(`&:focus-within`,Fe)]),z(`selected`,[Q(null,[I(`arrow`,`color: var(--n-arrow-color-active-hover);`),I(`icon`,`color: var(--n-item-icon-color-active-hover);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover);
 `,[D(`a`,`color: var(--n-item-text-color-active-hover);`),I(`extra`,`color: var(--n-item-text-color-active-hover);`)])])]),z(`child-active`,[Q(null,[I(`arrow`,`color: var(--n-arrow-color-child-active-hover);`),I(`icon`,`color: var(--n-item-icon-color-child-active-hover);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover);
 `,[D(`a`,`color: var(--n-item-text-color-child-active-hover);`),I(`extra`,`color: var(--n-item-text-color-child-active-hover);`)])])]),z(`selected`,[Q(null,[D(`&::before`,`background-color: var(--n-item-color-active-hover);`)])]),Q(null,Fe)]),I(`icon`,`
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
 `),I(`arrow`,`
 grid-area: arrow;
 font-size: 16px;
 color: var(--n-arrow-color);
 transform: rotate(180deg);
 opacity: 1;
 transition:
 color .3s var(--n-bezier),
 transform 0.2s var(--n-bezier),
 opacity 0.2s var(--n-bezier);
 `),w(`menu-item-content-header`,`
 grid-area: content;
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 opacity: 1;
 white-space: nowrap;
 color: var(--n-item-text-color);
 `,[D(`a`,`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[D(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),I(`extra`,`
 font-size: .93em;
 color: var(--n-group-text-color);
 transition: color .3s var(--n-bezier);
 `)])]),w(`submenu`,`
 cursor: pointer;
 position: relative;
 margin-top: 6px;
 `,[w(`menu-item-content`,`
 height: var(--n-item-height);
 `),w(`submenu-children`,`
 overflow: hidden;
 padding: 0;
 `,[de({duration:`.2s`})])]),w(`menu-item-group`,[w(`menu-item-group-title`,`
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
 `)])]),w(`menu-tooltip`,[D(`a`,`
 color: inherit;
 text-decoration: none;
 `)]),w(`menu-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function Q(e,t){return[z(`hover`,e,t),D(`&:hover`,e,t)]}var Re=g({name:`MenuDivider`,setup(){let{mergedClsPrefixRef:n,isHorizontalRef:r}=o(X);return()=>r.value?null:(M(),e(`div`,{key:1,class:t(`${n.value}-menu-divider`)},null,2))}}),ze=g({name:`ChevronDownFilled`,render(){return(()=>{let e=h(`f3af82a2aab086a5`);return e[0]||=c(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[c(`path`,{d:`M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z`,fill:`currentColor`})],-1)})()}}),Be=[`onClick`],Ve=g({name:`MenuOptionContent`,props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(e){let{props:t}=o(X);return{menuProps:t,style:i(()=>{let{paddingLeft:t}=e;return{paddingLeft:t&&`${t}px`}}),iconStyle:i(()=>{let{maxIconSize:t,activeIconSize:n,iconMarginRight:r}=e;return{width:`${t}px`,height:`${t}px`,fontSize:`${n}px`,marginRight:`${r}px`}})}},render(){let{clsPrefix:n,tmNode:r,menuProps:{renderIcon:i,renderLabel:a,renderExtra:o,expandIcon:s}}=this,u=i?i(r.rawNode):Y(this.icon);return(()=>{let i=h(`7bb10afc6caf8fa4`);return M(),e(`div`,{onClick:e=>{this.onClick?.(e)},role:`none`,class:t([`${n}-menu-item-content`,{[`${n}-menu-item-content--selected`]:this.selected,[`${n}-menu-item-content--collapsed`]:this.collapsed,[`${n}-menu-item-content--child-active`]:this.childActive,[`${n}-menu-item-content--disabled`]:this.disabled,[`${n}-menu-item-content--hover`]:this.hover}]),style:L(this.style)},[O(()=>u&&(M(),e(`div`,{class:t(`${n}-menu-item-content__icon`),style:L(this.iconStyle),role:`none`},[O(()=>[u])],6))),c(`div`,{class:t(`${n}-menu-item-content-header`),role:`none`},[this.isEllipsisPlaceholder?(M(),e(l,{key:0},[O(()=>this.title)],64)):(M(),e(l,{key:1},[a?(M(),e(l,{key:0},[O(()=>a(r.rawNode))],64)):(M(),e(l,{key:1},[O(()=>Y(this.title))],64))],64)),this.extra||o?(M(),e(`span`,{key:2,class:t(`${n}-menu-item-content-header__extra`)},[i[0]||=O(` `,-1),o?(M(),e(l,{key:0},[O(()=>o(r.rawNode))],64)):(M(),e(l,{key:1},[O(()=>Y(this.extra))],64))],2)):O(()=>null)],2),this.showArrow?(M(),f(S,{key:0,ariaHidden:!0,class:t(`${n}-menu-item-content__arrow`),clsPrefix:n},{default:()=>s?s(r.rawNode):(M(),f(ze,{key:1}))},1032,[`class`,`clsPrefix`])):O(()=>null)],14,Be)})()}}),He=8;function $(e){let t=o(X),{props:n,mergedCollapsedRef:r}=t,a=o(Pe,null),s=o(Z,null),c=i(()=>n.mode===`horizontal`),l=i(()=>c.value?n.dropdownPlacement:`tmNodes`in e?`right-start`:`right`),u=i(()=>Math.max(n.collapsedIconSize??n.iconSize,n.iconSize));return{dropdownPlacement:l,activeIconSize:i(()=>!c.value&&e.root&&r.value?n.collapsedIconSize??n.iconSize:n.iconSize),maxIconSize:u,paddingLeft:i(()=>{if(c.value)return;let{collapsedWidth:t,indent:i,rootIndent:o}=n,{root:l,isGroup:d}=e,f=o===void 0?i:o;return l?r.value?t/2-u.value/2:f:s&&typeof s.paddingLeftRef.value==`number`?r.value?t/2-u.value/2:i/2+s.paddingLeftRef.value:a&&typeof a.paddingLeftRef.value==`number`?(d?i/2:i)+a.paddingLeftRef.value:0}),iconMarginRight:i(()=>{let{collapsedWidth:t,indent:i,rootIndent:a}=n,{value:o}=u,{root:s}=e;return c.value||!s||!r.value?He:(a===void 0?i:a)+o+He-(t+o)/2}),NMenu:t,NSubmenu:a,NMenuOptionGroup:s}}var Ue={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},We={...Ue,tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function},Ge=H(We),Ke=g({name:`MenuOption`,props:We,setup(e){let t=$(e),{NSubmenu:n,NMenu:r,NMenuOptionGroup:a}=t,{props:o,mergedClsPrefixRef:s,mergedCollapsedRef:c}=r,l=n?n.mergedDisabledRef:a?a.mergedDisabledRef:{value:!1},u=i(()=>l.value||e.disabled);function d(t){let{onClick:n}=e;n&&n(t)}function f(t){u.value||(r.doSelect(e.internalKey,e.tmNode.rawNode),d(t))}return{mergedClsPrefix:s,dropdownPlacement:t.dropdownPlacement,paddingLeft:t.paddingLeft,iconMarginRight:t.iconMarginRight,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,mergedTheme:r.mergedThemeRef,menuProps:o,dropdownEnabled:x(()=>e.root&&c.value&&o.mode!==`horizontal`&&!u.value),selected:x(()=>r.mergedValueRef.value===e.internalKey),mergedDisabled:u,handleClick:f}},render(){let{mergedClsPrefix:t,mergedTheme:n,tmNode:r,menuProps:{renderLabel:i,nodeProps:o}}=this,s=o?.(r.rawNode);return M(),e(`div`,a(s,{role:`menuitem`,class:[`${t}-menu-item`,s?.class]}),[(M(),f(le,{theme:n.peers.Tooltip,themeOverrides:n.peerOverrides.Tooltip,trigger:`hover`,placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:[`menu-tooltip`]},{default:()=>i?i(r.rawNode):Y(this.title),trigger:()=>(M(),f(Ve,{tmNode:r,clsPrefix:t,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick},null,8,[`tmNode`,`clsPrefix`,`paddingLeft`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`selected`,`title`,`extra`,`disabled`,`icon`,`onClick`]))},1032,[`theme`,`themeOverrides`,`placement`,`disabled`]))],16)}}),qe={...Ue,tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}},Je=H(qe),Ye=g({name:`MenuOptionGroup`,props:qe,setup(n){let r=$(n),{NSubmenu:s}=r,u=i(()=>s?.mergedDisabledRef.value?!0:n.tmNode.disabled);P(Z,{paddingLeftRef:r.paddingLeft,mergedDisabledRef:u});let{mergedClsPrefixRef:d,props:f}=o(X);return function(){let{value:i}=d,o=r.paddingLeft.value,{nodeProps:s}=f,u=s?.(n.tmNode.rawNode);return(()=>{let r=h(`45eca6a63be5028b`);return M(),e(`div`,{class:t(`${i}-menu-item-group`),role:`group`},[c(`div`,a(u,{class:[`${i}-menu-item-group-title`,u?.class],style:[u?.style||``,o===void 0?``:`padding-left: ${o}px;`]}),[O(()=>Y(n.title)),n.extra?(M(),e(l,{key:0},[r[0]||=O(` `,-1),O(()=>Y(n.extra))],64)):O(()=>null)],16),c(`div`,null,[O(()=>n.tmNodes.map(e=>rt(e,f)))])],2)})()}}}),Xe=[`aria-expanded`,`id`],Ze=[`aria-expanded`,`id`],Qe={...Ue,rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean},$e=H(Qe),et=g({name:`Submenu`,props:Qe,setup(e){let t=$(e),{NMenu:n,NSubmenu:r}=t,{props:a,mergedCollapsedRef:o,mergedThemeRef:s}=n,c=i(()=>{let{disabled:t}=e;return r?.mergedDisabledRef.value||a.disabled?!0:t}),l=E(!1);P(Pe,{paddingLeftRef:t.paddingLeft,mergedDisabledRef:c}),P(Z,null);function u(){let{onClick:t}=e;t&&t()}function d(){c.value||(o.value||n.toggleExpand(e.internalKey),u())}function f(e){l.value=e}return{menuProps:a,mergedTheme:s,doSelect:n.doSelect,inverted:n.invertedRef,isHorizontal:n.isHorizontalRef,mergedClsPrefix:n.mergedClsPrefixRef,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,iconMarginRight:t.iconMarginRight,dropdownPlacement:t.dropdownPlacement,dropdownShow:l,paddingLeft:t.paddingLeft,mergedDisabled:c,mergedValue:n.mergedValueRef,childActive:x(()=>e.virtualChildActive??n.activePathRef.value.includes(e.internalKey)),collapsed:i(()=>a.mode===`horizontal`?!1:o.value?!0:!n.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:i(()=>!c.value&&(a.mode===`horizontal`||o.value)),handlePopoverShowChange:f,handleClick:d}},render(){let{mergedClsPrefix:n,menuProps:{renderIcon:r,renderLabel:i}}=this,o=()=>{let{isHorizontal:t,paddingLeft:n,collapsed:r,mergedDisabled:i,maxIconSize:o,activeIconSize:s,title:c,childActive:l,icon:u,handleClick:d,menuProps:{nodeProps:p},dropdownShow:m,iconMarginRight:h,tmNode:g,mergedClsPrefix:_,isEllipsisPlaceholder:v,extra:y}=this,b=p?.(g.rawNode);return M(),e(`div`,a(b,{class:[`${_}-menu-item`,b?.class],role:`menuitem`}),[(M(),f(Ve,{tmNode:g,paddingLeft:n,collapsed:r,disabled:i,iconMarginRight:h,maxIconSize:o,activeIconSize:s,title:c,extra:y,showArrow:!t,childActive:l,clsPrefix:_,icon:u,hover:m,onClick:d,isEllipsisPlaceholder:v},null,8,[`tmNode`,`paddingLeft`,`collapsed`,`disabled`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`title`,`extra`,`showArrow`,`childActive`,`clsPrefix`,`icon`,`hover`,`onClick`,`isEllipsisPlaceholder`]))],16)},s=()=>(M(),f(y,null,{default:()=>{let{tmNodes:r,collapsed:i}=this;return i?null:(M(),e(`div`,{key:1,class:t(`${n}-submenu-children`),role:`menu`},[O(()=>r.map(e=>rt(e,this.menuProps)))],2))}},1024));return this.root?(M(),f(ue,a({key:2,size:`large`,trigger:`hover`},this.menuProps?.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:`14px`,optionIconSizeLarge:`18px`},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:r,renderLabel:i}),{default:()=>(M(),e(`div`,{class:t(`${n}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[O(()=>o()),this.isHorizontal?O(()=>null):(M(),e(l,{key:1},[O(()=>s())],64))],10,Xe))},1040,[`themeOverrides`,`theme`,`value`,`disabled`,`placement`,`keyField`,`labelField`,`childrenField`,`onUpdateShow`,`options`,`onSelect`,`inverted`,`renderIcon`,`renderLabel`])):(M(),e(`div`,{key:3,class:t(`${n}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[O(()=>o()),O(()=>s())],10,Ze))}});function tt(e){return e.type===`divider`||e.type===`render`}function nt(e){return e.type===`divider`}function rt(e,t){let{rawNode:r}=e,{show:i}=r;if(i===!1)return null;if(tt(r))return nt(r)?(M(),f(Re,a({key:e.key},r.props),null,16)):null;let{labelField:o}=t,{key:s,level:c,isGroup:l}=e,u={...r,title:r.title||r[o],extra:r.titleExtra||r.extra,key:s,internalKey:s,level:c,root:c===0,isGroup:l};return e.children?e.isGroup?n(Ye,q(u,Je,{tmNode:e,tmNodes:e.children,key:s})):n(et,q(u,$e,{key:s,rawNodes:r[t.childrenField],tmNodes:e.children,tmNode:e})):n(Ke,q(u,Ge,{key:s,tmNode:e}))}var it={...A.props,options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},disabledField:{type:String,default:`disabled`},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:`vertical`},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:`bottom`},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array},at=g({name:`Menu`,inheritAttrs:!1,props:it,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=d(e),r=A(`Menu`,`-menu`,Le,me,e,t),a=o(ve,null),s=i(()=>{let{collapsed:t}=e;if(t!==void 0)return t;if(a){let{collapseModeRef:e,collapsedRef:t}=a;if(e.value===`width`)return t.value??!1}return!1}),c=i(()=>{let{keyField:t,childrenField:n,disabledField:r}=e;return B(e.items||e.options,{getIgnored(e){return tt(e)},getChildren(e){return e[n]},getDisabled(e){return e[r]},getKey(e){return e[t]??e.name}})}),l=i(()=>new Set(c.value.treeNodes.map(e=>e.key))),{watchProps:u}=e,p=E(null);u?.includes(`defaultValue`)?b(()=>{p.value=e.defaultValue}):p.value=e.defaultValue;let m=F(e,`value`),h=J(m,p),g=E([]),_=()=>{g.value=e.defaultExpandAll?c.value.getNonLeafKeys():e.defaultExpandedNames||e.defaultExpandedKeys||c.value.getPath(h.value,{includeSelf:!1}).keyPath};u?.includes(`defaultExpandedKeys`)?b(_):_();let v=ie(e,[`expandedNames`,`expandedKeys`]),y=J(v,g),x=i(()=>c.value.treeNodes),S=i(()=>c.value.getPath(h.value).keyPath);P(X,{props:e,mergedCollapsedRef:s,mergedThemeRef:r,mergedValueRef:h,mergedExpandedKeysRef:y,activePathRef:S,mergedClsPrefixRef:t,isHorizontalRef:i(()=>e.mode===`horizontal`),invertedRef:F(e,`inverted`),doSelect:C,toggleExpand:D});function C(t,n){let{"onUpdate:value":r,onUpdateValue:i,onSelect:a}=e;i&&U(i,t,n),r&&U(r,t,n),a&&U(a,t,n),p.value=t}function w(t){let{"onUpdate:expandedKeys":n,onUpdateExpandedKeys:r,onExpandedNamesChange:i,onOpenNamesChange:a}=e;n&&U(n,t),r&&U(r,t),i&&U(i,t),a&&U(a,t),g.value=t}function D(t){let n=Array.from(y.value),r=n.findIndex(e=>e===t);if(~r)n.splice(r,1);else{if(e.accordion&&l.value.has(t)){let e=n.findIndex(e=>l.value.has(e));e>-1&&n.splice(e,1)}n.push(t)}w(n)}let O=t=>{let n=c.value.getPath(t??h.value,{includeSelf:!1}).keyPath;if(!n.length)return;let r=Array.from(y.value),i=new Set([...r,...n]);e.accordion&&l.value.forEach(e=>{i.has(e)&&!n.includes(e)&&i.delete(e)}),w(Array.from(i))},k=i(()=>{let{inverted:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{borderRadius:a,borderColorHorizontal:o,fontSize:s,itemHeight:c,dividerColor:l}=i,u={"--n-divider-color":l,"--n-bezier":n,"--n-font-size":s,"--n-border-color-horizontal":o,"--n-border-radius":a,"--n-item-height":c};return t?(u[`--n-group-text-color`]=i.groupTextColorInverted,u[`--n-color`]=i.colorInverted,u[`--n-item-text-color`]=i.itemTextColorInverted,u[`--n-item-text-color-hover`]=i.itemTextColorHoverInverted,u[`--n-item-text-color-active`]=i.itemTextColorActiveInverted,u[`--n-item-text-color-child-active`]=i.itemTextColorChildActiveInverted,u[`--n-item-text-color-child-active-hover`]=i.itemTextColorChildActiveInverted,u[`--n-item-text-color-active-hover`]=i.itemTextColorActiveHoverInverted,u[`--n-item-icon-color`]=i.itemIconColorInverted,u[`--n-item-icon-color-hover`]=i.itemIconColorHoverInverted,u[`--n-item-icon-color-active`]=i.itemIconColorActiveInverted,u[`--n-item-icon-color-active-hover`]=i.itemIconColorActiveHoverInverted,u[`--n-item-icon-color-child-active`]=i.itemIconColorChildActiveInverted,u[`--n-item-icon-color-child-active-hover`]=i.itemIconColorChildActiveHoverInverted,u[`--n-item-icon-color-collapsed`]=i.itemIconColorCollapsedInverted,u[`--n-item-text-color-horizontal`]=i.itemTextColorHorizontalInverted,u[`--n-item-text-color-hover-horizontal`]=i.itemTextColorHoverHorizontalInverted,u[`--n-item-text-color-active-horizontal`]=i.itemTextColorActiveHorizontalInverted,u[`--n-item-text-color-child-active-horizontal`]=i.itemTextColorChildActiveHorizontalInverted,u[`--n-item-text-color-child-active-hover-horizontal`]=i.itemTextColorChildActiveHoverHorizontalInverted,u[`--n-item-text-color-active-hover-horizontal`]=i.itemTextColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-horizontal`]=i.itemIconColorHorizontalInverted,u[`--n-item-icon-color-hover-horizontal`]=i.itemIconColorHoverHorizontalInverted,u[`--n-item-icon-color-active-horizontal`]=i.itemIconColorActiveHorizontalInverted,u[`--n-item-icon-color-active-hover-horizontal`]=i.itemIconColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-child-active-horizontal`]=i.itemIconColorChildActiveHorizontalInverted,u[`--n-item-icon-color-child-active-hover-horizontal`]=i.itemIconColorChildActiveHoverHorizontalInverted,u[`--n-arrow-color`]=i.arrowColorInverted,u[`--n-arrow-color-hover`]=i.arrowColorHoverInverted,u[`--n-arrow-color-active`]=i.arrowColorActiveInverted,u[`--n-arrow-color-active-hover`]=i.arrowColorActiveHoverInverted,u[`--n-arrow-color-child-active`]=i.arrowColorChildActiveInverted,u[`--n-arrow-color-child-active-hover`]=i.arrowColorChildActiveHoverInverted,u[`--n-item-color-hover`]=i.itemColorHoverInverted,u[`--n-item-color-active`]=i.itemColorActiveInverted,u[`--n-item-color-active-hover`]=i.itemColorActiveHoverInverted,u[`--n-item-color-active-collapsed`]=i.itemColorActiveCollapsedInverted):(u[`--n-group-text-color`]=i.groupTextColor,u[`--n-color`]=i.color,u[`--n-item-text-color`]=i.itemTextColor,u[`--n-item-text-color-hover`]=i.itemTextColorHover,u[`--n-item-text-color-active`]=i.itemTextColorActive,u[`--n-item-text-color-child-active`]=i.itemTextColorChildActive,u[`--n-item-text-color-child-active-hover`]=i.itemTextColorChildActiveHover,u[`--n-item-text-color-active-hover`]=i.itemTextColorActiveHover,u[`--n-item-icon-color`]=i.itemIconColor,u[`--n-item-icon-color-hover`]=i.itemIconColorHover,u[`--n-item-icon-color-active`]=i.itemIconColorActive,u[`--n-item-icon-color-active-hover`]=i.itemIconColorActiveHover,u[`--n-item-icon-color-child-active`]=i.itemIconColorChildActive,u[`--n-item-icon-color-child-active-hover`]=i.itemIconColorChildActiveHover,u[`--n-item-icon-color-collapsed`]=i.itemIconColorCollapsed,u[`--n-item-text-color-horizontal`]=i.itemTextColorHorizontal,u[`--n-item-text-color-hover-horizontal`]=i.itemTextColorHoverHorizontal,u[`--n-item-text-color-active-horizontal`]=i.itemTextColorActiveHorizontal,u[`--n-item-text-color-child-active-horizontal`]=i.itemTextColorChildActiveHorizontal,u[`--n-item-text-color-child-active-hover-horizontal`]=i.itemTextColorChildActiveHoverHorizontal,u[`--n-item-text-color-active-hover-horizontal`]=i.itemTextColorActiveHoverHorizontal,u[`--n-item-icon-color-horizontal`]=i.itemIconColorHorizontal,u[`--n-item-icon-color-hover-horizontal`]=i.itemIconColorHoverHorizontal,u[`--n-item-icon-color-active-horizontal`]=i.itemIconColorActiveHorizontal,u[`--n-item-icon-color-active-hover-horizontal`]=i.itemIconColorActiveHoverHorizontal,u[`--n-item-icon-color-child-active-horizontal`]=i.itemIconColorChildActiveHorizontal,u[`--n-item-icon-color-child-active-hover-horizontal`]=i.itemIconColorChildActiveHoverHorizontal,u[`--n-arrow-color`]=i.arrowColor,u[`--n-arrow-color-hover`]=i.arrowColorHover,u[`--n-arrow-color-active`]=i.arrowColorActive,u[`--n-arrow-color-active-hover`]=i.arrowColorActiveHover,u[`--n-arrow-color-child-active`]=i.arrowColorChildActive,u[`--n-arrow-color-child-active-hover`]=i.arrowColorChildActiveHover,u[`--n-item-color-hover`]=i.itemColorHover,u[`--n-item-color-active`]=i.itemColorActive,u[`--n-item-color-active-hover`]=i.itemColorActiveHover,u[`--n-item-color-active-collapsed`]=i.itemColorActiveCollapsed),u}),j=n?T(`menu`,i(()=>e.inverted?`a`:`b`),k,e):void 0,N=ne(),ee=E(null),te=E(null),I=!0,L=()=>{I?I=!1:ee.value?.sync({showAllItemsBeforeCalculate:!0})};function R(){return document.getElementById(N)}let z=E(-1);function re(t){z.value=e.options.length-t}function V(e){e||(z.value=-1)}let H=i(()=>{let t=z.value;return{children:t===-1?[]:e.options.slice(t)}}),W=i(()=>{let{childrenField:t,disabledField:n,keyField:r}=e;return B([H.value],{getIgnored(e){return tt(e)},getChildren(e){return e[t]},getDisabled(e){return e[n]},getKey(e){return e[r]??e.name}})}),G=i(()=>B([{}]).treeNodes[0]);function K(){if(z.value===-1)return M(),f(et,{root:!0,level:0,key:`__ellpisisGroupPlaceholder__`,internalKey:`__ellpisisGroupPlaceholder__`,title:`···`,tmNode:G.value,domId:N,isEllipsisPlaceholder:!0},null,8,[`tmNode`,`domId`]);let e=W.value.treeNodes[0],t=S.value,n=!!e.children?.some(e=>t.includes(e.key));return M(),f(et,{level:0,root:!0,key:`__ellpisisGroup__`,internalKey:`__ellpisisGroup__`,title:`···`,virtualChildActive:n,tmNode:e,domId:N,rawNodes:e.rawNode.children||[],tmNodes:e.children||[],isEllipsisPlaceholder:!0},null,8,[`virtualChildActive`,`tmNode`,`domId`,`rawNodes`,`tmNodes`])}return{mergedClsPrefix:t,controlledExpandedKeys:v,uncontrolledExpanededKeys:g,mergedExpandedKeys:y,uncontrolledValue:p,mergedValue:h,activePath:S,tmNodes:x,mergedTheme:r,mergedCollapsed:s,cssVars:n?void 0:k,themeClass:j?.themeClass,overflowRef:ee,counterRef:te,updateCounter:()=>{},onResize:L,onUpdateOverflow:V,onUpdateCount:re,renderCounter:K,getCounter:R,onRender:j?.onRender,showOption:O,deriveResponsiveState:L}},render(){let{mergedClsPrefix:e,mode:t,themeClass:r,onRender:i}=this;i?.();let o=()=>this.tmNodes.map(e=>rt(e,this.$props)),s=t===`horizontal`&&this.responsive,c=()=>n(`div`,a(this.$attrs,{role:t===`horizontal`?`menubar`:`menu`,class:[`${e}-menu`,r,`${e}-menu--${t}`,s&&`${e}-menu--responsive`,this.mergedCollapsed&&`${e}-menu--collapsed`],style:this.cssVars}),s?(M(),f(re,{key:2,ref:`overflowRef`,onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:o,counter:this.renderCounter},1032,[`onUpdateOverflow`,`getCounter`,`onUpdateCount`,`updateCounter`])):o());return s?(M(),f(se,{key:3,onResize:this.onResize},{default:c},1032,[`onResize`])):c()}}),ot={class:`sider-footer`},st=oe(g({__name:`LayoutView`,setup(e){let{t,locale:n}=j(),a=pe(),o=i(()=>[{label:t(`nav.devices`),key:`/devices`},{label:t(`nav.expect`),key:`/expect`},{label:t(`nav.history`),key:`/history`},{label:t(`nav.settings`),key:`/settings`,disabled:!0}]);function s(e){a.push(e)}function l(){n.value=n.value===`zh`?`en`:`zh`}function u(){ae(),a.push(`/login`)}return(e,n)=>{let i=N(`router-view`);return M(),f(R(we),{position:`absolute`,class:`root`},{default:k(()=>[r(R(Ne),{bordered:``,width:200,class:`sider`},{default:k(()=>[n[0]||=c(`div`,{class:`brand`},`EmbedFlow`,-1),r(R(at),{options:o.value,"onUpdate:value":s},null,8,[`options`]),c(`div`,ot,[r(R(K),{quaternary:``,size:`small`,onClick:R(fe)},{default:k(()=>[_(m(R(he).isDark?`Light`:`Dark`),1)]),_:1},8,[`onClick`]),r(R(K),{quaternary:``,size:`small`,onClick:l},{default:k(()=>[_(m(R(t)(`nav.language`)),1)]),_:1}),r(R(K),{quaternary:``,size:`small`,onClick:u},{default:k(()=>[_(m(R(t)(`nav.logout`)),1)]),_:1})])]),_:1}),r(R(Te),{"content-style":`padding: 16px;`,class:`content`},{default:k(()=>[r(i)]),_:1})]),_:1})}}}),[[`__scopeId`,`data-v-6bff097f`]]);export{st as default};