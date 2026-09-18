import{Cn as e,E as t,Ht as n,Kt as r,L as i,Lt as a,Nt as o,O as s,P as c,Rt as l,S as u,Sn as d,T as f,Ut as p,Vt as m,Wt as h,Yt as g,bn as _,bt as v,c as y,d as b,f as x,g as S,gn as C,gt as w,h as T,ht as E,in as D,k as O,ln as ee,m as k,n as A,p as te,qt as j,rn as M,sn as ne,un as N,v as P,vt as F,y as re,yn as I,yt as L,zt as R}from"./vue-i18n-pPSPCP6m.js";import{o as ie,p as ae,t as z}from"./create-_LNlS2dw.js";import{t as B}from"./format-length-Bb0m8EKo.js";import{I as V,O as H,S as U,g as W,l as G,n as oe,t as se,v as ce}from"./_plugin-vue_export-helper-cUzcR9Kp.js";import{c as K,f as q}from"./fade-in-scale-up.cssr-BU07VhSR.js";import{i as le,n as ue,t as de}from"./Dropdown-DfzOVuB2.js";import{D as fe,E as J,i as pe,n as me,o as he,r as ge}from"./index-BHHyBIUJ.js";function _e(e){let{baseColor:t,textColor2:n,bodyColor:r,cardColor:i,dividerColor:a,actionColor:o,scrollbarColor:s,scrollbarColorHover:c,invertedColor:l}=e;return{textColor:n,textColorInverted:`#FFF`,color:r,colorEmbedded:o,headerColor:i,headerColorInverted:l,footerColor:o,footerColorInverted:l,headerBorderColor:a,headerBorderColorInverted:l,footerBorderColor:a,footerBorderColorInverted:l,siderBorderColor:a,siderBorderColorInverted:l,siderColor:i,siderColorInverted:l,siderToggleButtonBorder:`1px solid ${a}`,siderToggleButtonColor:t,siderToggleButtonIconColor:n,siderToggleButtonIconColorInverted:n,siderToggleBarColor:u(r,s),siderToggleBarColorHover:u(r,c),__invertScrollbar:`true`}}var ve=te({name:`Layout`,common:P,peers:{Scrollbar:S},self:_e}),ye=i(`n-layout-sider`),be={type:String,default:`static`},xe=w(`layout`,`
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
 `),L(`absolute-positioned`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),Se={embedded:Boolean,position:be,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:``},hasSider:Boolean,siderPlacement:{type:String,default:`left`}},Ce=i(`n-layout`);function we(e){return h({name:e?`LayoutContent`:`Layout`,props:{...k.props,...Se},setup(e){let t=C(null),n=C(null),{mergedClsPrefixRef:r,inlineThemeDisabled:i}=c(e),o=k(`Layout`,`-layout`,xe,ve,e,r);function s(r,i){if(e.nativeScrollbar){let{value:e}=t;e&&(i===void 0?e.scrollTo(r):e.scrollTo(r,i))}else{let{value:e}=n;e&&e.scrollTo(r,i)}}D(Ce,e);let l=0,u=0,d=t=>{let n=t.target;l=n.scrollLeft,u=n.scrollTop,e.onScroll?.(t)};U(()=>{if(e.nativeScrollbar){let e=t.value;e&&(e.scrollTop=u,e.scrollLeft=l)}});let f={display:`flex`,flexWrap:`nowrap`,width:`100%`,flexDirection:`row`},p={scrollTo:s},m=a(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=o.value;return{"--n-bezier":t,"--n-color":e.embedded?n.colorEmbedded:n.color,"--n-text-color":n.textColor}}),h=i?T(`layout`,a(()=>e.embedded?`e`:``),m,e):void 0;return{mergedClsPrefix:r,scrollableElRef:t,scrollbarInstRef:n,hasSiderStyle:f,mergedTheme:o,handleNativeElScroll:d,cssVars:i?void 0:m,themeClass:h?.themeClass,onRender:h?.onRender,...p}},render(){let{mergedClsPrefix:n,hasSider:r}=this;this.onRender?.();let i=r?this.hasSiderStyle:void 0,a=[this.themeClass,e&&`${n}-layout-content`,`${n}-layout`,`${n}-layout--${this.position}-positioned`];return M(),m(`div`,{class:t(a),style:d(this.cssVars)},[this.nativeScrollbar?(M(),m(`div`,{key:0,ref:`scrollableElRef`,class:t([`${n}-layout-scroll-container`,this.contentClass]),style:d([this.contentStyle,i]),onScroll:this.handleNativeElScroll},[O(()=>this.$slots.default?.())],46,[`onScroll`])):(M(),R(W,g({key:1},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,i]}),s(this.$slots),1040,[`onScroll`,`theme`,`themeOverrides`,`contentClass`,`contentStyle`]))],6)}})}var Te=we(!1),Ee=we(!0),De=w(`layout-sider`,`
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
`,[L(`bordered`,[F(`border`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),F(`left-placement`,[L(`bordered`,[F(`border`,`
 right: 0;
 `)])]),L(`right-placement`,`
 justify-content: flex-start;
 `,[L(`bordered`,[F(`border`,`
 left: 0;
 `)]),L(`collapsed`,[w(`layout-toggle-button`,[w(`base-icon`,`
 transform: rotate(180deg);
 `)]),w(`layout-toggle-bar`,[E(`&:hover`,[F(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),F(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])])]),w(`layout-toggle-button`,`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[w(`base-icon`,`
 transform: rotate(0);
 `)]),w(`layout-toggle-bar`,`
 left: -28px;
 transform: rotate(180deg);
 `,[E(`&:hover`,[F(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),F(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})])])]),L(`collapsed`,[w(`layout-toggle-bar`,[E(`&:hover`,[F(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),F(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])]),w(`layout-toggle-button`,[w(`base-icon`,`
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
 `,[F(`top, bottom`,`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),F(`bottom`,`
 position: absolute;
 top: 34px;
 `),E(`&:hover`,[F(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),F(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})]),F(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color)`}),E(`&:hover`,[F(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color-hover)`})])]),F(`border`,`
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
 `),L(`show-content`,[w(`layout-sider-scroll-container`,{opacity:1})]),L(`absolute-positioned`,`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]),Oe=[`onClick`],ke=h({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:e}=this;return M(),m(`div`,{onClick:this.onClick,class:t(`${e}-layout-toggle-bar`)},[l(`div`,{class:t(`${e}-layout-toggle-bar__top`)},null,2),l(`div`,{class:t(`${e}-layout-toggle-bar__bottom`)},null,2)],10,Oe)}}),Ae=[`onClick`],je=h({name:`LayoutToggleButton`,props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:e}=this;return M(),m(`div`,{class:t(`${e}-layout-toggle-button`),onClick:this.onClick},[(M(),R(x,{clsPrefix:e},{default:()=>(M(),R(le))},1032,[`clsPrefix`]))],10,Ae)}}),Me=[`onTransitionend`],Ne={position:be,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:``},collapseMode:{type:String,default:`transform`},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function},Pe=h({name:`LayoutSider`,props:{...k.props,...Ne},setup(e){let t=j(Ce),n=C(null),r=C(null),i=C(e.defaultCollapsed),o=q(I(e,`collapsed`),i),s=a(()=>B(o.value?e.collapsedWidth:e.width)),l=a(()=>e.collapseMode===`transform`?{minWidth:B(e.width)}:{}),u=a(()=>t?t.siderPlacement:`left`);function d(t,i){if(e.nativeScrollbar){let{value:e}=n;e&&(i===void 0?e.scrollTo(t):e.scrollTo(t,i))}else{let{value:e}=r;e&&e.scrollTo(t,i)}}function f(){let{"onUpdate:collapsed":t,onUpdateCollapsed:n,onExpand:r,onCollapse:a}=e,{value:s}=o;n&&H(n,!s),t&&H(t,!s),i.value=!s,s?r&&H(r):a&&H(a)}let p=0,m=0,h=t=>{let n=t.target;p=n.scrollLeft,m=n.scrollTop,e.onScroll?.(t)};U(()=>{if(e.nativeScrollbar){let e=n.value;e&&(e.scrollTop=m,e.scrollLeft=p)}}),D(ye,{collapsedRef:o,collapseModeRef:I(e,`collapseMode`)});let{mergedClsPrefixRef:g,inlineThemeDisabled:_}=c(e),v=k(`Layout`,`-layout-sider`,De,ve,e,g);function y(t){t.propertyName===`max-width`&&(o.value?e.onAfterLeave?.():e.onAfterEnter?.())}let b={scrollTo:d},x=a(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=v.value,{siderToggleButtonColor:r,siderToggleButtonBorder:i,siderToggleBarColor:a,siderToggleBarColorHover:o}=n,s={"--n-bezier":t,"--n-toggle-button-color":r,"--n-toggle-button-border":i,"--n-toggle-bar-color":a,"--n-toggle-bar-color-hover":o};return e.inverted?(s[`--n-color`]=n.siderColorInverted,s[`--n-text-color`]=n.textColorInverted,s[`--n-border-color`]=n.siderBorderColorInverted,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColorInverted,s.__invertScrollbar=n.__invertScrollbar):(s[`--n-color`]=n.siderColor,s[`--n-text-color`]=n.textColor,s[`--n-border-color`]=n.siderBorderColor,s[`--n-toggle-button-icon-color`]=n.siderToggleButtonIconColor),s}),S=_?T(`layout-sider`,a(()=>e.inverted?`a`:`b`),x,e):void 0;return{scrollableElRef:n,scrollbarInstRef:r,mergedClsPrefix:g,mergedTheme:v,styleMaxWidth:s,mergedCollapsed:o,scrollContainerStyle:l,siderPlacement:u,handleNativeElScroll:h,handleTransitionend:y,handleTriggerClick:f,inlineThemeDisabled:_,cssVars:x,themeClass:S?.themeClass,onRender:S?.onRender,...b}},render(){let{mergedClsPrefix:e,mergedCollapsed:n,showTrigger:r}=this;return this.onRender?.(),M(),m(`aside`,{class:t([`${e}-layout-sider`,this.themeClass,`${e}-layout-sider--${this.position}-positioned`,`${e}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${e}-layout-sider--bordered`,n&&`${e}-layout-sider--collapsed`,(!n||this.showCollapsedContent)&&`${e}-layout-sider--show-content`]),onTransitionend:this.handleTransitionend,style:d([this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:B(this.width)}])},[this.nativeScrollbar?(M(),m(`div`,{key:1,class:t([`${e}-layout-sider-scroll-container`,this.contentClass]),onScroll:this.handleNativeElScroll,style:d([this.scrollContainerStyle,{overflow:`auto`},this.contentStyle]),ref:`scrollableElRef`},[O(()=>this.$slots.default?.())],46,[`onScroll`])):(M(),R(W,g({key:0},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar===`true`?{colorHover:`rgba(255, 255, 255, .4)`,color:`rgba(255, 255, 255, .3)`}:void 0}),s(this.$slots),1040,[`onScroll`,`style`,`contentStyle`,`contentClass`,`theme`,`themeOverrides`,`builtinThemeOverrides`])),r?(M(),m(o,{key:2},[r===`bar`?(M(),R(ke,{key:0,clsPrefix:e,class:t(n?this.collapsedTriggerClass:this.triggerClass),style:d(n?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`])):(M(),R(je,{key:1,clsPrefix:e,class:t(n?this.collapsedTriggerClass:this.triggerClass),style:d(n?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`]))],64)):O(()=>null),this.bordered?(M(),m(`div`,{key:4,class:t(`${e}-layout-sider__border`)},null,2)):O(()=>null)],46,Me)}}),Y=i(`n-menu`),Fe=i(`n-submenu`),X=i(`n-menu-item-group`),Ie=[E(`&::before`,`background-color: var(--n-item-color-hover);`),F(`arrow`,`
 color: var(--n-arrow-color-hover);
 `),F(`icon`,`
 color: var(--n-item-icon-color-hover);
 `),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover);
 `,[E(`a`,`
 color: var(--n-item-text-color-hover);
 `),F(`extra`,`
 color: var(--n-item-text-color-hover);
 `)])],Le=[F(`icon`,`
 color: var(--n-item-icon-color-hover-horizontal);
 `),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover-horizontal);
 `,[E(`a`,`
 color: var(--n-item-text-color-hover-horizontal);
 `),F(`extra`,`
 color: var(--n-item-text-color-hover-horizontal);
 `)])],Re=E([w(`menu`,`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[L(`horizontal`,`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[w(`submenu`,`margin: 0;`),w(`menu-item`,`margin: 0;`),w(`menu-item-content`,`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[E(`&::before`,`display: none;`),L(`selected`,`border-bottom: 2px solid var(--n-border-color-horizontal)`)]),w(`menu-item-content`,[L(`selected`,[F(`icon`,`color: var(--n-item-icon-color-active-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-horizontal);
 `,[E(`a`,`color: var(--n-item-text-color-active-horizontal);`),F(`extra`,`color: var(--n-item-text-color-active-horizontal);`)])]),L(`child-active`,`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[E(`a`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `),F(`extra`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),F(`icon`,`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),v(`disabled`,[v(`selected, child-active`,[E(`&:focus-within`,Le)]),L(`selected`,[Z(null,[F(`icon`,`color: var(--n-item-icon-color-active-hover-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[E(`a`,`color: var(--n-item-text-color-active-hover-horizontal);`),F(`extra`,`color: var(--n-item-text-color-active-hover-horizontal);`)])])]),L(`child-active`,[Z(null,[F(`icon`,`color: var(--n-item-icon-color-child-active-hover-horizontal);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[E(`a`,`color: var(--n-item-text-color-child-active-hover-horizontal);`),F(`extra`,`color: var(--n-item-text-color-child-active-hover-horizontal);`)])])]),Z(`border-bottom: 2px solid var(--n-border-color-horizontal);`,Le)]),w(`menu-item-content-header`,[E(`a`,`color: var(--n-item-text-color-horizontal);`)])])]),v(`responsive`,[w(`menu-item-content-header`,`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),L(`collapsed`,[w(`menu-item-content`,[L(`selected`,[E(`&::before`,`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),w(`menu-item-content-header`,`opacity: 0;`),F(`arrow`,`opacity: 0;`),F(`icon`,`color: var(--n-item-icon-color-collapsed);`)])]),w(`menu-item`,`
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
 `,[E(`> *`,`z-index: 1;`),E(`&::before`,`
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
 `),L(`disabled`,`
 opacity: .45;
 cursor: not-allowed;
 `),L(`collapsed`,[F(`arrow`,`transform: rotate(0);`)]),L(`selected`,[E(`&::before`,`background-color: var(--n-item-color-active);`),F(`arrow`,`color: var(--n-arrow-color-active);`),F(`icon`,`color: var(--n-item-icon-color-active);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active);
 `,[E(`a`,`color: var(--n-item-text-color-active);`),F(`extra`,`color: var(--n-item-text-color-active);`)])]),L(`child-active`,[w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active);
 `,[E(`a`,`
 color: var(--n-item-text-color-child-active);
 `),F(`extra`,`
 color: var(--n-item-text-color-child-active);
 `)]),F(`arrow`,`
 color: var(--n-arrow-color-child-active);
 `),F(`icon`,`
 color: var(--n-item-icon-color-child-active);
 `)]),v(`disabled`,[v(`selected, child-active`,[E(`&:focus-within`,Ie)]),L(`selected`,[Z(null,[F(`arrow`,`color: var(--n-arrow-color-active-hover);`),F(`icon`,`color: var(--n-item-icon-color-active-hover);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover);
 `,[E(`a`,`color: var(--n-item-text-color-active-hover);`),F(`extra`,`color: var(--n-item-text-color-active-hover);`)])])]),L(`child-active`,[Z(null,[F(`arrow`,`color: var(--n-arrow-color-child-active-hover);`),F(`icon`,`color: var(--n-item-icon-color-child-active-hover);`),w(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover);
 `,[E(`a`,`color: var(--n-item-text-color-child-active-hover);`),F(`extra`,`color: var(--n-item-text-color-child-active-hover);`)])])]),L(`selected`,[Z(null,[E(`&::before`,`background-color: var(--n-item-color-active-hover);`)])]),Z(null,Ie)]),F(`icon`,`
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
 `),F(`arrow`,`
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
 `,[E(`a`,`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[E(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),F(`extra`,`
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
 `,[fe({duration:`.2s`})])]),w(`menu-item-group`,[w(`menu-item-group-title`,`
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
 `)])]),w(`menu-tooltip`,[E(`a`,`
 color: inherit;
 text-decoration: none;
 `)]),w(`menu-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function Z(e,t){return[L(`hover`,e,t),E(`&:hover`,e,t)]}var ze=h({name:`MenuDivider`,setup(){let{mergedClsPrefixRef:e,isHorizontalRef:n}=j(Y);return()=>n.value?null:(M(),m(`div`,{key:1,class:t(`${e.value}-menu-divider`)},null,2))}}),Be=h({name:`ChevronDownFilled`,render(){return(()=>{let e=f(`f3af82a2aab086a5`);return e[0]||=l(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[l(`path`,{d:`M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z`,fill:`currentColor`})],-1)})()}}),Ve=[`onClick`],He=h({name:`MenuOptionContent`,props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(e){let{props:t}=j(Y);return{menuProps:t,style:a(()=>{let{paddingLeft:t}=e;return{paddingLeft:t&&`${t}px`}}),iconStyle:a(()=>{let{maxIconSize:t,activeIconSize:n,iconMarginRight:r}=e;return{width:`${t}px`,height:`${t}px`,fontSize:`${n}px`,marginRight:`${r}px`}})}},render(){let{clsPrefix:e,tmNode:n,menuProps:{renderIcon:r,renderLabel:i,renderExtra:a,expandIcon:s}}=this,c=r?r(n.rawNode):J(this.icon);return(()=>{let r=f(`7bb10afc6caf8fa4`);return M(),m(`div`,{onClick:e=>{this.onClick?.(e)},role:`none`,class:t([`${e}-menu-item-content`,{[`${e}-menu-item-content--selected`]:this.selected,[`${e}-menu-item-content--collapsed`]:this.collapsed,[`${e}-menu-item-content--child-active`]:this.childActive,[`${e}-menu-item-content--disabled`]:this.disabled,[`${e}-menu-item-content--hover`]:this.hover}]),style:d(this.style)},[O(()=>c&&(M(),m(`div`,{class:t(`${e}-menu-item-content__icon`),style:d(this.iconStyle),role:`none`},[O(()=>[c])],6))),l(`div`,{class:t(`${e}-menu-item-content-header`),role:`none`},[this.isEllipsisPlaceholder?(M(),m(o,{key:0},[O(()=>this.title)],64)):(M(),m(o,{key:1},[i?(M(),m(o,{key:0},[O(()=>i(n.rawNode))],64)):(M(),m(o,{key:1},[O(()=>J(this.title))],64))],64)),this.extra||a?(M(),m(`span`,{key:2,class:t(`${e}-menu-item-content-header__extra`)},[r[0]||=O(` `,-1),a?(M(),m(o,{key:0},[O(()=>a(n.rawNode))],64)):(M(),m(o,{key:1},[O(()=>J(this.extra))],64))],2)):O(()=>null)],2),this.showArrow?(M(),R(x,{key:0,ariaHidden:!0,class:t(`${e}-menu-item-content__arrow`),clsPrefix:e},{default:()=>s?s(n.rawNode):(M(),R(Be,{key:1}))},1032,[`class`,`clsPrefix`])):O(()=>null)],14,Ve)})()}}),Ue=8;function Q(e){let t=j(Y),{props:n,mergedCollapsedRef:r}=t,i=j(Fe,null),o=j(X,null),s=a(()=>n.mode===`horizontal`),c=a(()=>s.value?n.dropdownPlacement:`tmNodes`in e?`right-start`:`right`),l=a(()=>Math.max(n.collapsedIconSize??n.iconSize,n.iconSize));return{dropdownPlacement:c,activeIconSize:a(()=>!s.value&&e.root&&r.value?n.collapsedIconSize??n.iconSize:n.iconSize),maxIconSize:l,paddingLeft:a(()=>{if(s.value)return;let{collapsedWidth:t,indent:a,rootIndent:c}=n,{root:u,isGroup:d}=e,f=c===void 0?a:c;return u?r.value?t/2-l.value/2:f:o&&typeof o.paddingLeftRef.value==`number`?r.value?t/2-l.value/2:a/2+o.paddingLeftRef.value:i&&typeof i.paddingLeftRef.value==`number`?(d?a/2:a)+i.paddingLeftRef.value:0}),iconMarginRight:a(()=>{let{collapsedWidth:t,indent:i,rootIndent:a}=n,{value:o}=l,{root:c}=e;return s.value||!c||!r.value?Ue:(a===void 0?i:a)+o+Ue-(t+o)/2}),NMenu:t,NSubmenu:i,NMenuOptionGroup:o}}var $={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},We={...$,tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function},Ge=V(We),Ke=h({name:`MenuOption`,props:We,setup(e){let t=Q(e),{NSubmenu:n,NMenu:r,NMenuOptionGroup:i}=t,{props:o,mergedClsPrefixRef:s,mergedCollapsedRef:c}=r,l=n?n.mergedDisabledRef:i?i.mergedDisabledRef:{value:!1},u=a(()=>l.value||e.disabled);function d(t){let{onClick:n}=e;n&&n(t)}function f(t){u.value||(r.doSelect(e.internalKey,e.tmNode.rawNode),d(t))}return{mergedClsPrefix:s,dropdownPlacement:t.dropdownPlacement,paddingLeft:t.paddingLeft,iconMarginRight:t.iconMarginRight,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,mergedTheme:r.mergedThemeRef,menuProps:o,dropdownEnabled:b(()=>e.root&&c.value&&o.mode!==`horizontal`&&!u.value),selected:b(()=>r.mergedValueRef.value===e.internalKey),mergedDisabled:u,handleClick:f}},render(){let{mergedClsPrefix:e,mergedTheme:t,tmNode:n,menuProps:{renderLabel:r,nodeProps:i}}=this,a=i?.(n.rawNode);return M(),m(`div`,g(a,{role:`menuitem`,class:[`${e}-menu-item`,a?.class]}),[(M(),R(ue,{theme:t.peers.Tooltip,themeOverrides:t.peerOverrides.Tooltip,trigger:`hover`,placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:[`menu-tooltip`]},{default:()=>r?r(n.rawNode):J(this.title),trigger:()=>(M(),R(He,{tmNode:n,clsPrefix:e,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick},null,8,[`tmNode`,`clsPrefix`,`paddingLeft`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`selected`,`title`,`extra`,`disabled`,`icon`,`onClick`]))},1032,[`theme`,`themeOverrides`,`placement`,`disabled`]))],16)}}),qe={...$,tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}},Je=V(qe),Ye=h({name:`MenuOptionGroup`,props:qe,setup(e){let n=Q(e),{NSubmenu:r}=n,i=a(()=>r?.mergedDisabledRef.value?!0:e.tmNode.disabled);D(X,{paddingLeftRef:n.paddingLeft,mergedDisabledRef:i});let{mergedClsPrefixRef:s,props:c}=j(Y);return function(){let{value:r}=s,i=n.paddingLeft.value,{nodeProps:a}=c,u=a?.(e.tmNode.rawNode);return(()=>{let n=f(`45eca6a63be5028b`);return M(),m(`div`,{class:t(`${r}-menu-item-group`),role:`group`},[l(`div`,g(u,{class:[`${r}-menu-item-group-title`,u?.class],style:[u?.style||``,i===void 0?``:`padding-left: ${i}px;`]}),[O(()=>J(e.title)),e.extra?(M(),m(o,{key:0},[n[0]||=O(` `,-1),O(()=>J(e.extra))],64)):O(()=>null)],16),l(`div`,null,[O(()=>e.tmNodes.map(e=>rt(e,c)))])],2)})()}}}),Xe=[`aria-expanded`,`id`],Ze=[`aria-expanded`,`id`],Qe={...$,rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean},$e=V(Qe),et=h({name:`Submenu`,props:Qe,setup(e){let t=Q(e),{NMenu:n,NSubmenu:r}=t,{props:i,mergedCollapsedRef:o,mergedThemeRef:s}=n,c=a(()=>{let{disabled:t}=e;return r?.mergedDisabledRef.value||i.disabled?!0:t}),l=C(!1);D(Fe,{paddingLeftRef:t.paddingLeft,mergedDisabledRef:c}),D(X,null);function u(){let{onClick:t}=e;t&&t()}function d(){c.value||(o.value||n.toggleExpand(e.internalKey),u())}function f(e){l.value=e}return{menuProps:i,mergedTheme:s,doSelect:n.doSelect,inverted:n.invertedRef,isHorizontal:n.isHorizontalRef,mergedClsPrefix:n.mergedClsPrefixRef,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,iconMarginRight:t.iconMarginRight,dropdownPlacement:t.dropdownPlacement,dropdownShow:l,paddingLeft:t.paddingLeft,mergedDisabled:c,mergedValue:n.mergedValueRef,childActive:b(()=>e.virtualChildActive??n.activePathRef.value.includes(e.internalKey)),collapsed:a(()=>i.mode===`horizontal`?!1:o.value?!0:!n.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:a(()=>!c.value&&(i.mode===`horizontal`||o.value)),handlePopoverShowChange:f,handleClick:d}},render(){let{mergedClsPrefix:e,menuProps:{renderIcon:n,renderLabel:r}}=this,i=()=>{let{isHorizontal:e,paddingLeft:t,collapsed:n,mergedDisabled:r,maxIconSize:i,activeIconSize:a,title:o,childActive:s,icon:c,handleClick:l,menuProps:{nodeProps:u},dropdownShow:d,iconMarginRight:f,tmNode:p,mergedClsPrefix:h,isEllipsisPlaceholder:_,extra:v}=this,y=u?.(p.rawNode);return M(),m(`div`,g(y,{class:[`${h}-menu-item`,y?.class],role:`menuitem`}),[(M(),R(He,{tmNode:p,paddingLeft:t,collapsed:n,disabled:r,iconMarginRight:f,maxIconSize:i,activeIconSize:a,title:o,extra:v,showArrow:!e,childActive:s,clsPrefix:h,icon:c,hover:d,onClick:l,isEllipsisPlaceholder:_},null,8,[`tmNode`,`paddingLeft`,`collapsed`,`disabled`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`title`,`extra`,`showArrow`,`childActive`,`clsPrefix`,`icon`,`hover`,`onClick`,`isEllipsisPlaceholder`]))],16)},a=()=>(M(),R(y,null,{default:()=>{let{tmNodes:n,collapsed:r}=this;return r?null:(M(),m(`div`,{key:1,class:t(`${e}-submenu-children`),role:`menu`},[O(()=>n.map(e=>rt(e,this.menuProps)))],2))}},1024));return this.root?(M(),R(de,g({key:2,size:`large`,trigger:`hover`},this.menuProps?.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:`14px`,optionIconSizeLarge:`18px`},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:n,renderLabel:r}),{default:()=>(M(),m(`div`,{class:t(`${e}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[O(()=>i()),this.isHorizontal?O(()=>null):(M(),m(o,{key:1},[O(()=>a())],64))],10,Xe))},1040,[`themeOverrides`,`theme`,`value`,`disabled`,`placement`,`keyField`,`labelField`,`childrenField`,`onUpdateShow`,`options`,`onSelect`,`inverted`,`renderIcon`,`renderLabel`])):(M(),m(`div`,{key:3,class:t(`${e}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[O(()=>i()),O(()=>a())],10,Ze))}});function tt(e){return e.type===`divider`||e.type===`render`}function nt(e){return e.type===`divider`}function rt(e,t){let{rawNode:n}=e,{show:i}=n;if(i===!1)return null;if(tt(n))return nt(n)?(M(),R(ze,g({key:e.key},n.props),null,16)):null;let{labelField:a}=t,{key:o,level:s,isGroup:c}=e,l={...n,title:n.title||n[a],extra:n.titleExtra||n.extra,key:o,internalKey:o,level:s,root:s===0,isGroup:c};return e.children?e.isGroup?r(Ye,K(l,Je,{tmNode:e,tmNodes:e.children,key:o})):r(et,K(l,$e,{key:o,rawNodes:n[t.childrenField],tmNodes:e.children,tmNode:e})):r(Ke,K(l,Ge,{key:o,tmNode:e}))}var it={...k.props,options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},disabledField:{type:String,default:`disabled`},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:`vertical`},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:`bottom`},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array},at=h({name:`Menu`,inheritAttrs:!1,props:it,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n}=c(e),r=k(`Menu`,`-menu`,Re,he,e,t),i=j(ye,null),o=a(()=>{let{collapsed:t}=e;if(t!==void 0)return t;if(i){let{collapseModeRef:e,collapsedRef:t}=i;if(e.value===`width`)return t.value??!1}return!1}),s=a(()=>{let{keyField:t,childrenField:n,disabledField:r}=e;return z(e.items||e.options,{getIgnored(e){return tt(e)},getChildren(e){return e[n]},getDisabled(e){return e[r]},getKey(e){return e[t]??e.name}})}),l=a(()=>new Set(s.value.treeNodes.map(e=>e.key))),{watchProps:u}=e,d=C(null);u?.includes(`defaultValue`)?ee(()=>{d.value=e.defaultValue}):d.value=e.defaultValue;let f=I(e,`value`),p=q(f,d),m=C([]),h=()=>{m.value=e.defaultExpandAll?s.value.getNonLeafKeys():e.defaultExpandedNames||e.defaultExpandedKeys||s.value.getPath(p.value,{includeSelf:!1}).keyPath};u?.includes(`defaultExpandedKeys`)?ee(h):h();let g=ae(e,[`expandedNames`,`expandedKeys`]),_=q(g,m),v=a(()=>s.value.treeNodes),y=a(()=>s.value.getPath(p.value).keyPath);D(Y,{props:e,mergedCollapsedRef:o,mergedThemeRef:r,mergedValueRef:p,mergedExpandedKeysRef:_,activePathRef:y,mergedClsPrefixRef:t,isHorizontalRef:a(()=>e.mode===`horizontal`),invertedRef:I(e,`inverted`),doSelect:b,toggleExpand:S});function b(t,n){let{"onUpdate:value":r,onUpdateValue:i,onSelect:a}=e;i&&H(i,t,n),r&&H(r,t,n),a&&H(a,t,n),d.value=t}function x(t){let{"onUpdate:expandedKeys":n,onUpdateExpandedKeys:r,onExpandedNamesChange:i,onOpenNamesChange:a}=e;n&&H(n,t),r&&H(r,t),i&&H(i,t),a&&H(a,t),m.value=t}function S(t){let n=Array.from(_.value),r=n.findIndex(e=>e===t);if(~r)n.splice(r,1);else{if(e.accordion&&l.value.has(t)){let e=n.findIndex(e=>l.value.has(e));e>-1&&n.splice(e,1)}n.push(t)}x(n)}let w=t=>{let n=s.value.getPath(t??p.value,{includeSelf:!1}).keyPath;if(!n.length)return;let r=Array.from(_.value),i=new Set([...r,...n]);e.accordion&&l.value.forEach(e=>{i.has(e)&&!n.includes(e)&&i.delete(e)}),x(Array.from(i))},E=a(()=>{let{inverted:t}=e,{common:{cubicBezierEaseInOut:n},self:i}=r.value,{borderRadius:a,borderColorHorizontal:o,fontSize:s,itemHeight:c,dividerColor:l}=i,u={"--n-divider-color":l,"--n-bezier":n,"--n-font-size":s,"--n-border-color-horizontal":o,"--n-border-radius":a,"--n-item-height":c};return t?(u[`--n-group-text-color`]=i.groupTextColorInverted,u[`--n-color`]=i.colorInverted,u[`--n-item-text-color`]=i.itemTextColorInverted,u[`--n-item-text-color-hover`]=i.itemTextColorHoverInverted,u[`--n-item-text-color-active`]=i.itemTextColorActiveInverted,u[`--n-item-text-color-child-active`]=i.itemTextColorChildActiveInverted,u[`--n-item-text-color-child-active-hover`]=i.itemTextColorChildActiveInverted,u[`--n-item-text-color-active-hover`]=i.itemTextColorActiveHoverInverted,u[`--n-item-icon-color`]=i.itemIconColorInverted,u[`--n-item-icon-color-hover`]=i.itemIconColorHoverInverted,u[`--n-item-icon-color-active`]=i.itemIconColorActiveInverted,u[`--n-item-icon-color-active-hover`]=i.itemIconColorActiveHoverInverted,u[`--n-item-icon-color-child-active`]=i.itemIconColorChildActiveInverted,u[`--n-item-icon-color-child-active-hover`]=i.itemIconColorChildActiveHoverInverted,u[`--n-item-icon-color-collapsed`]=i.itemIconColorCollapsedInverted,u[`--n-item-text-color-horizontal`]=i.itemTextColorHorizontalInverted,u[`--n-item-text-color-hover-horizontal`]=i.itemTextColorHoverHorizontalInverted,u[`--n-item-text-color-active-horizontal`]=i.itemTextColorActiveHorizontalInverted,u[`--n-item-text-color-child-active-horizontal`]=i.itemTextColorChildActiveHorizontalInverted,u[`--n-item-text-color-child-active-hover-horizontal`]=i.itemTextColorChildActiveHoverHorizontalInverted,u[`--n-item-text-color-active-hover-horizontal`]=i.itemTextColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-horizontal`]=i.itemIconColorHorizontalInverted,u[`--n-item-icon-color-hover-horizontal`]=i.itemIconColorHoverHorizontalInverted,u[`--n-item-icon-color-active-horizontal`]=i.itemIconColorActiveHorizontalInverted,u[`--n-item-icon-color-active-hover-horizontal`]=i.itemIconColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-child-active-horizontal`]=i.itemIconColorChildActiveHorizontalInverted,u[`--n-item-icon-color-child-active-hover-horizontal`]=i.itemIconColorChildActiveHoverHorizontalInverted,u[`--n-arrow-color`]=i.arrowColorInverted,u[`--n-arrow-color-hover`]=i.arrowColorHoverInverted,u[`--n-arrow-color-active`]=i.arrowColorActiveInverted,u[`--n-arrow-color-active-hover`]=i.arrowColorActiveHoverInverted,u[`--n-arrow-color-child-active`]=i.arrowColorChildActiveInverted,u[`--n-arrow-color-child-active-hover`]=i.arrowColorChildActiveHoverInverted,u[`--n-item-color-hover`]=i.itemColorHoverInverted,u[`--n-item-color-active`]=i.itemColorActiveInverted,u[`--n-item-color-active-hover`]=i.itemColorActiveHoverInverted,u[`--n-item-color-active-collapsed`]=i.itemColorActiveCollapsedInverted):(u[`--n-group-text-color`]=i.groupTextColor,u[`--n-color`]=i.color,u[`--n-item-text-color`]=i.itemTextColor,u[`--n-item-text-color-hover`]=i.itemTextColorHover,u[`--n-item-text-color-active`]=i.itemTextColorActive,u[`--n-item-text-color-child-active`]=i.itemTextColorChildActive,u[`--n-item-text-color-child-active-hover`]=i.itemTextColorChildActiveHover,u[`--n-item-text-color-active-hover`]=i.itemTextColorActiveHover,u[`--n-item-icon-color`]=i.itemIconColor,u[`--n-item-icon-color-hover`]=i.itemIconColorHover,u[`--n-item-icon-color-active`]=i.itemIconColorActive,u[`--n-item-icon-color-active-hover`]=i.itemIconColorActiveHover,u[`--n-item-icon-color-child-active`]=i.itemIconColorChildActive,u[`--n-item-icon-color-child-active-hover`]=i.itemIconColorChildActiveHover,u[`--n-item-icon-color-collapsed`]=i.itemIconColorCollapsed,u[`--n-item-text-color-horizontal`]=i.itemTextColorHorizontal,u[`--n-item-text-color-hover-horizontal`]=i.itemTextColorHoverHorizontal,u[`--n-item-text-color-active-horizontal`]=i.itemTextColorActiveHorizontal,u[`--n-item-text-color-child-active-horizontal`]=i.itemTextColorChildActiveHorizontal,u[`--n-item-text-color-child-active-hover-horizontal`]=i.itemTextColorChildActiveHoverHorizontal,u[`--n-item-text-color-active-hover-horizontal`]=i.itemTextColorActiveHoverHorizontal,u[`--n-item-icon-color-horizontal`]=i.itemIconColorHorizontal,u[`--n-item-icon-color-hover-horizontal`]=i.itemIconColorHoverHorizontal,u[`--n-item-icon-color-active-horizontal`]=i.itemIconColorActiveHorizontal,u[`--n-item-icon-color-active-hover-horizontal`]=i.itemIconColorActiveHoverHorizontal,u[`--n-item-icon-color-child-active-horizontal`]=i.itemIconColorChildActiveHorizontal,u[`--n-item-icon-color-child-active-hover-horizontal`]=i.itemIconColorChildActiveHoverHorizontal,u[`--n-arrow-color`]=i.arrowColor,u[`--n-arrow-color-hover`]=i.arrowColorHover,u[`--n-arrow-color-active`]=i.arrowColorActive,u[`--n-arrow-color-active-hover`]=i.arrowColorActiveHover,u[`--n-arrow-color-child-active`]=i.arrowColorChildActive,u[`--n-arrow-color-child-active-hover`]=i.arrowColorChildActiveHover,u[`--n-item-color-hover`]=i.itemColorHover,u[`--n-item-color-active`]=i.itemColorActive,u[`--n-item-color-active-hover`]=i.itemColorActiveHover,u[`--n-item-color-active-collapsed`]=i.itemColorActiveCollapsed),u}),O=n?T(`menu`,a(()=>e.inverted?`a`:`b`),E,e):void 0,A=re(),te=C(null),ne=C(null),N=!0,P=()=>{N?N=!1:te.value?.sync({showAllItemsBeforeCalculate:!0})};function F(){return document.getElementById(A)}let L=C(-1);function ie(t){L.value=e.options.length-t}function B(e){e||(L.value=-1)}let V=a(()=>{let t=L.value;return{children:t===-1?[]:e.options.slice(t)}}),U=a(()=>{let{childrenField:t,disabledField:n,keyField:r}=e;return z([V.value],{getIgnored(e){return tt(e)},getChildren(e){return e[t]},getDisabled(e){return e[n]},getKey(e){return e[r]??e.name}})}),W=a(()=>z([{}]).treeNodes[0]);function G(){if(L.value===-1)return M(),R(et,{root:!0,level:0,key:`__ellpisisGroupPlaceholder__`,internalKey:`__ellpisisGroupPlaceholder__`,title:`···`,tmNode:W.value,domId:A,isEllipsisPlaceholder:!0},null,8,[`tmNode`,`domId`]);let e=U.value.treeNodes[0],t=y.value,n=!!e.children?.some(e=>t.includes(e.key));return M(),R(et,{level:0,root:!0,key:`__ellpisisGroup__`,internalKey:`__ellpisisGroup__`,title:`···`,virtualChildActive:n,tmNode:e,domId:A,rawNodes:e.rawNode.children||[],tmNodes:e.children||[],isEllipsisPlaceholder:!0},null,8,[`virtualChildActive`,`tmNode`,`domId`,`rawNodes`,`tmNodes`])}return{mergedClsPrefix:t,controlledExpandedKeys:g,uncontrolledExpanededKeys:m,mergedExpandedKeys:_,uncontrolledValue:d,mergedValue:p,activePath:y,tmNodes:v,mergedTheme:r,mergedCollapsed:o,cssVars:n?void 0:E,themeClass:O?.themeClass,overflowRef:te,counterRef:ne,updateCounter:()=>{},onResize:P,onUpdateOverflow:B,onUpdateCount:ie,renderCounter:G,getCounter:F,onRender:O?.onRender,showOption:w,deriveResponsiveState:P}},render(){let{mergedClsPrefix:e,mode:t,themeClass:n,onRender:i}=this;i?.();let a=()=>this.tmNodes.map(e=>rt(e,this.$props)),o=t===`horizontal`&&this.responsive,s=()=>r(`div`,g(this.$attrs,{role:t===`horizontal`?`menubar`:`menu`,class:[`${e}-menu`,n,`${e}-menu--${t}`,o&&`${e}-menu--responsive`,this.mergedCollapsed&&`${e}-menu--collapsed`],style:this.cssVars}),o?(M(),R(ie,{key:2,ref:`overflowRef`,onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:a,counter:this.renderCounter},1032,[`onUpdateOverflow`,`getCounter`,`onUpdateCount`,`updateCounter`])):a());return o?(M(),R(ce,{key:3,onResize:this.onResize},{default:s},1032,[`onResize`])):s()}}),ot={class:`sider-footer`},st=se(h({__name:`LayoutView`,setup(t){let{t:r,locale:i}=A(),o=me(),s=a(()=>[{label:r(`nav.devices`),key:`/devices`},{label:r(`nav.expect`),key:`/expect`},{label:r(`nav.history`),key:`/history`},{label:r(`nav.settings`),key:`/settings`,disabled:!0}]);function c(e){o.push(e)}function u(){i.value=i.value===`zh`?`en`:`zh`}function d(){oe(),o.push(`/login`)}return(t,i)=>{let a=ne(`router-view`);return M(),R(_(Te),{position:`absolute`,class:`root`},{default:N(()=>[p(_(Pe),{bordered:``,width:200,class:`sider`},{default:N(()=>[i[0]||=l(`div`,{class:`brand`},`EmbedFlow`,-1),p(_(at),{options:s.value,"onUpdate:value":c},null,8,[`options`]),l(`div`,ot,[p(_(G),{quaternary:``,size:`small`,onClick:_(pe)},{default:N(()=>[n(e(_(ge).isDark?`Light`:`Dark`),1)]),_:1},8,[`onClick`]),p(_(G),{quaternary:``,size:`small`,onClick:u},{default:N(()=>[n(e(_(r)(`nav.language`)),1)]),_:1}),p(_(G),{quaternary:``,size:`small`,onClick:d},{default:N(()=>[n(e(_(r)(`nav.logout`)),1)]),_:1})])]),_:1}),p(_(Ee),{"content-style":`padding: 16px;`,class:`content`},{default:N(()=>[p(a)]),_:1})]),_:1})}}}),[[`__scopeId`,`data-v-6bff097f`]]);export{st as default};