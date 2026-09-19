import{C as e,F as t,G as n,P as r,R as i,S as a,T as o,U as s,W as c,X as l,_ as u,_t as d,d as f,et as p,f as m,h,j as g,n as _,o as v,p as y,tt as b,v as x,vt as S,y as C,yt as w}from"./vue-i18n-CU1juWHN.js";import{At as T,E,L as D,O,Ot as ee,P as te,S as k,T as A,bt as j,d as M,f as N,g as P,gt as F,h as I,ht as L,k as R,m as z,o as ne,p as re,v as ie,vt as B,y as ae,yt as V}from"./light-uQ0rL05w.js";import{m as H,v as U}from"./event-ww-MKoYV.js";import{t as W}from"./render-LoiffuQC.js";import{h as oe,i as G,l as se,t as ce}from"./Select-DRR-liTq.js";import{t as K}from"./format-length-IO90KiIs.js";import{C as le,F as q,T as ue,U as J,g as Y,i as de,k as fe}from"./http-CYJUkLIH.js";import{i as pe,n as me,t as he}from"./Dropdown-1wB8O5WC.js";import{D as ge,i as _e,n as ve,o as ye,r as be}from"./index-DSu5_6pc.js";import{t as xe}from"./project-DEOGGql3.js";import{t as Se}from"./_plugin-vue_export-helper-BDNMzG2s.js";function Ce(e){let{baseColor:t,textColor2:n,bodyColor:r,cardColor:i,dividerColor:a,actionColor:o,scrollbarColor:s,scrollbarColorHover:c,invertedColor:l}=e;return{textColor:n,textColorInverted:`#FFF`,color:r,colorEmbedded:o,headerColor:i,headerColorInverted:l,footerColor:o,footerColorInverted:l,headerBorderColor:a,headerBorderColorInverted:l,footerBorderColor:a,footerBorderColorInverted:l,siderBorderColor:a,siderBorderColorInverted:l,siderColor:i,siderColorInverted:l,siderToggleButtonBorder:`1px solid ${a}`,siderToggleButtonColor:t,siderToggleButtonIconColor:n,siderToggleButtonIconColorInverted:n,siderToggleBarColor:k(r,s),siderToggleBarColorHover:k(r,c),__invertScrollbar:`true`}}var we=re({name:`Layout`,common:ie,peers:{Scrollbar:P},self:Ce}),Te=D(`n-layout-sider`),Ee={type:String,default:`static`},De=F(`layout`,`
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
`,[F(`layout-scroll-container`,`
 overflow-x: hidden;
 box-sizing: border-box;
 height: 100%;
 `),V(`absolute-positioned`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),Oe={embedded:Boolean,position:Ee,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:``},hasSider:Boolean,siderPlacement:{type:String,default:`left`}},ke=D(`n-layout`);function Ae(e){return C({name:e?`LayoutContent`:`Layout`,props:{...z.props,...Oe},setup(e){let n=l(null),r=l(null),{mergedClsPrefixRef:i,inlineThemeDisabled:a}=te(e),o=z(`Layout`,`-layout`,De,we,e,i);function s(t,i){if(e.nativeScrollbar){let{value:e}=n;e&&(i===void 0?e.scrollTo(t):e.scrollTo(t,i))}else{let{value:e}=r;e&&e.scrollTo(t,i)}}t(ke,e);let c=0,u=0,d=t=>{let n=t.target;c=n.scrollLeft,u=n.scrollTop,e.onScroll?.(t)};fe(()=>{if(e.nativeScrollbar){let e=n.value;e&&(e.scrollTop=u,e.scrollLeft=c)}});let p={display:`flex`,flexWrap:`nowrap`,width:`100%`,flexDirection:`row`},m={scrollTo:s},h=f(()=>{let{common:{cubicBezierEaseInOut:t},self:n}=o.value;return{"--n-bezier":t,"--n-color":e.embedded?n.colorEmbedded:n.color,"--n-text-color":n.textColor}}),g=a?I(`layout`,f(()=>e.embedded?`e`:``),h,e):void 0;return{mergedClsPrefix:i,scrollableElRef:n,scrollbarInstRef:r,hasSiderStyle:p,mergedTheme:o,handleNativeElScroll:d,cssVars:a?void 0:h,themeClass:g?.themeClass,onRender:g?.onRender,...m}},render(){let{mergedClsPrefix:t,hasSider:n}=this;this.onRender?.();let i=n?this.hasSiderStyle:void 0,a=[this.themeClass,e&&`${t}-layout-content`,`${t}-layout`,`${t}-layout--${this.position}-positioned`];return r(),h(`div`,{class:E(a),style:S(this.cssVars)},[this.nativeScrollbar?(r(),h(`div`,{key:0,ref:`scrollableElRef`,class:E([`${t}-layout-scroll-container`,this.contentClass]),style:S([this.contentStyle,i]),onScroll:this.handleNativeElScroll},[R(()=>this.$slots.default?.())],46,[`onScroll`])):(r(),y(le,o({key:1},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,i]}),O(this.$slots),1040,[`onScroll`,`theme`,`themeOverrides`,`contentClass`,`contentStyle`]))],6)}})}var je=Ae(!1),Me=Ae(!0),Ne=F(`layout-sider`,`
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
`,[V(`bordered`,[B(`border`,`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),B(`left-placement`,[V(`bordered`,[B(`border`,`
 right: 0;
 `)])]),V(`right-placement`,`
 justify-content: flex-start;
 `,[V(`bordered`,[B(`border`,`
 left: 0;
 `)]),V(`collapsed`,[F(`layout-toggle-button`,[F(`base-icon`,`
 transform: rotate(180deg);
 `)]),F(`layout-toggle-bar`,[L(`&:hover`,[B(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),B(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])])]),F(`layout-toggle-button`,`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[F(`base-icon`,`
 transform: rotate(0);
 `)]),F(`layout-toggle-bar`,`
 left: -28px;
 transform: rotate(180deg);
 `,[L(`&:hover`,[B(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),B(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})])])]),V(`collapsed`,[F(`layout-toggle-bar`,[L(`&:hover`,[B(`top`,{transform:`rotate(-12deg) scale(1.15) translateY(-2px)`}),B(`bottom`,{transform:`rotate(12deg) scale(1.15) translateY(2px)`})])]),F(`layout-toggle-button`,[F(`base-icon`,`
 transform: rotate(0);
 `)])]),F(`layout-toggle-button`,`
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
 `,[F(`base-icon`,`
 transition: transform .3s var(--n-bezier);
 transform: rotate(180deg);
 `)]),F(`layout-toggle-bar`,`
 cursor: pointer;
 height: 72px;
 width: 32px;
 position: absolute;
 top: calc(50% - 36px);
 right: -28px;
 `,[B(`top, bottom`,`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),B(`bottom`,`
 position: absolute;
 top: 34px;
 `),L(`&:hover`,[B(`top`,{transform:`rotate(12deg) scale(1.15) translateY(-2px)`}),B(`bottom`,{transform:`rotate(-12deg) scale(1.15) translateY(2px)`})]),B(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color)`}),L(`&:hover`,[B(`top, bottom`,{backgroundColor:`var(--n-toggle-bar-color-hover)`})])]),B(`border`,`
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 width: 1px;
 transition: background-color .3s var(--n-bezier);
 `),F(`layout-sider-scroll-container`,`
 flex-grow: 1;
 flex-shrink: 0;
 box-sizing: border-box;
 height: 100%;
 opacity: 0;
 transition: opacity .3s var(--n-bezier);
 max-width: 100%;
 `),V(`show-content`,[F(`layout-sider-scroll-container`,{opacity:1})]),V(`absolute-positioned`,`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]),Pe=[`onClick`],Fe=C({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:e}=this;return r(),h(`div`,{onClick:this.onClick,class:E(`${e}-layout-toggle-bar`)},[m(`div`,{class:E(`${e}-layout-toggle-bar__top`)},null,2),m(`div`,{class:E(`${e}-layout-toggle-bar__bottom`)},null,2)],10,Pe)}}),Ie=[`onClick`],Le=C({name:`LayoutToggleButton`,props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){let{clsPrefix:e}=this;return r(),h(`div`,{class:E(`${e}-layout-toggle-button`),onClick:this.onClick},[(r(),y(N,{clsPrefix:e},{default:()=>(r(),y(pe))},1032,[`clsPrefix`]))],10,Ie)}}),Re=[`onTransitionend`],ze={position:Ee,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:``},collapseMode:{type:String,default:`transform`},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function},Be=C({name:`LayoutSider`,props:{...z.props,...ze},setup(n){let r=e(ke),i=l(null),a=l(null),o=l(n.defaultCollapsed),s=U(p(n,`collapsed`),o),c=f(()=>K(s.value?n.collapsedWidth:n.width)),u=f(()=>n.collapseMode===`transform`?{minWidth:K(n.width)}:{}),d=f(()=>r?r.siderPlacement:`left`);function m(e,t){if(n.nativeScrollbar){let{value:n}=i;n&&(t===void 0?n.scrollTo(e):n.scrollTo(e,t))}else{let{value:n}=a;n&&n.scrollTo(e,t)}}function h(){let{"onUpdate:collapsed":e,onUpdateCollapsed:t,onExpand:r,onCollapse:i}=n,{value:a}=s;t&&q(t,!a),e&&q(e,!a),o.value=!a,a?r&&q(r):i&&q(i)}let g=0,_=0,v=e=>{let t=e.target;g=t.scrollLeft,_=t.scrollTop,n.onScroll?.(e)};fe(()=>{if(n.nativeScrollbar){let e=i.value;e&&(e.scrollTop=_,e.scrollLeft=g)}}),t(Te,{collapsedRef:s,collapseModeRef:p(n,`collapseMode`)});let{mergedClsPrefixRef:y,inlineThemeDisabled:b}=te(n),x=z(`Layout`,`-layout-sider`,Ne,we,n,y);function S(e){e.propertyName===`max-width`&&(s.value?n.onAfterLeave?.():n.onAfterEnter?.())}let C={scrollTo:m},w=f(()=>{let{common:{cubicBezierEaseInOut:e},self:t}=x.value,{siderToggleButtonColor:r,siderToggleButtonBorder:i,siderToggleBarColor:a,siderToggleBarColorHover:o}=t,s={"--n-bezier":e,"--n-toggle-button-color":r,"--n-toggle-button-border":i,"--n-toggle-bar-color":a,"--n-toggle-bar-color-hover":o};return n.inverted?(s[`--n-color`]=t.siderColorInverted,s[`--n-text-color`]=t.textColorInverted,s[`--n-border-color`]=t.siderBorderColorInverted,s[`--n-toggle-button-icon-color`]=t.siderToggleButtonIconColorInverted,s.__invertScrollbar=t.__invertScrollbar):(s[`--n-color`]=t.siderColor,s[`--n-text-color`]=t.textColor,s[`--n-border-color`]=t.siderBorderColor,s[`--n-toggle-button-icon-color`]=t.siderToggleButtonIconColor),s}),T=b?I(`layout-sider`,f(()=>n.inverted?`a`:`b`),w,n):void 0;return{scrollableElRef:i,scrollbarInstRef:a,mergedClsPrefix:y,mergedTheme:x,styleMaxWidth:c,mergedCollapsed:s,scrollContainerStyle:u,siderPlacement:d,handleNativeElScroll:v,handleTransitionend:S,handleTriggerClick:h,inlineThemeDisabled:b,cssVars:w,themeClass:T?.themeClass,onRender:T?.onRender,...C}},render(){let{mergedClsPrefix:e,mergedCollapsed:t,showTrigger:n}=this;return this.onRender?.(),r(),h(`aside`,{class:E([`${e}-layout-sider`,this.themeClass,`${e}-layout-sider--${this.position}-positioned`,`${e}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${e}-layout-sider--bordered`,t&&`${e}-layout-sider--collapsed`,(!t||this.showCollapsedContent)&&`${e}-layout-sider--show-content`]),onTransitionend:this.handleTransitionend,style:S([this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:K(this.width)}])},[this.nativeScrollbar?(r(),h(`div`,{key:1,class:E([`${e}-layout-sider-scroll-container`,this.contentClass]),onScroll:this.handleNativeElScroll,style:S([this.scrollContainerStyle,{overflow:`auto`},this.contentStyle]),ref:`scrollableElRef`},[R(()=>this.$slots.default?.())],46,[`onScroll`])):(r(),y(le,o({key:0},this.scrollbarProps,{onScroll:this.onScroll,ref:`scrollbarInstRef`,style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar===`true`?{colorHover:`rgba(255, 255, 255, .4)`,color:`rgba(255, 255, 255, .3)`}:void 0}),O(this.$slots),1040,[`onScroll`,`style`,`contentStyle`,`contentClass`,`theme`,`themeOverrides`,`builtinThemeOverrides`])),n?(r(),h(v,{key:2},[n===`bar`?(r(),y(Fe,{key:0,clsPrefix:e,class:E(t?this.collapsedTriggerClass:this.triggerClass),style:S(t?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`])):(r(),y(Le,{key:1,clsPrefix:e,class:E(t?this.collapsedTriggerClass:this.triggerClass),style:S(t?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,[`clsPrefix`,`class`,`style`,`onClick`]))],64)):R(()=>null),this.bordered?(r(),h(`div`,{key:4,class:E(`${e}-layout-sider__border`)},null,2)):R(()=>null)],46,Re)}}),X=D(`n-menu`),Ve=D(`n-submenu`),Z=D(`n-menu-item-group`),He=[L(`&::before`,`background-color: var(--n-item-color-hover);`),B(`arrow`,`
 color: var(--n-arrow-color-hover);
 `),B(`icon`,`
 color: var(--n-item-icon-color-hover);
 `),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover);
 `,[L(`a`,`
 color: var(--n-item-text-color-hover);
 `),B(`extra`,`
 color: var(--n-item-text-color-hover);
 `)])],Ue=[B(`icon`,`
 color: var(--n-item-icon-color-hover-horizontal);
 `),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-hover-horizontal);
 `,[L(`a`,`
 color: var(--n-item-text-color-hover-horizontal);
 `),B(`extra`,`
 color: var(--n-item-text-color-hover-horizontal);
 `)])],We=L([F(`menu`,`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[V(`horizontal`,`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[F(`submenu`,`margin: 0;`),F(`menu-item`,`margin: 0;`),F(`menu-item-content`,`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[L(`&::before`,`display: none;`),V(`selected`,`border-bottom: 2px solid var(--n-border-color-horizontal)`)]),F(`menu-item-content`,[V(`selected`,[B(`icon`,`color: var(--n-item-icon-color-active-horizontal);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-horizontal);
 `,[L(`a`,`color: var(--n-item-text-color-active-horizontal);`),B(`extra`,`color: var(--n-item-text-color-active-horizontal);`)])]),V(`child-active`,`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[F(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[L(`a`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `),B(`extra`,`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),B(`icon`,`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),j(`disabled`,[j(`selected, child-active`,[L(`&:focus-within`,Ue)]),V(`selected`,[Q(null,[B(`icon`,`color: var(--n-item-icon-color-active-hover-horizontal);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[L(`a`,`color: var(--n-item-text-color-active-hover-horizontal);`),B(`extra`,`color: var(--n-item-text-color-active-hover-horizontal);`)])])]),V(`child-active`,[Q(null,[B(`icon`,`color: var(--n-item-icon-color-child-active-hover-horizontal);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[L(`a`,`color: var(--n-item-text-color-child-active-hover-horizontal);`),B(`extra`,`color: var(--n-item-text-color-child-active-hover-horizontal);`)])])]),Q(`border-bottom: 2px solid var(--n-border-color-horizontal);`,Ue)]),F(`menu-item-content-header`,[L(`a`,`color: var(--n-item-text-color-horizontal);`)])])]),j(`responsive`,[F(`menu-item-content-header`,`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),V(`collapsed`,[F(`menu-item-content`,[V(`selected`,[L(`&::before`,`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),F(`menu-item-content-header`,`opacity: 0;`),B(`arrow`,`opacity: 0;`),B(`icon`,`color: var(--n-item-icon-color-collapsed);`)])]),F(`menu-item`,`
 height: var(--n-item-height);
 margin-top: 6px;
 position: relative;
 `),F(`menu-item-content`,`
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
 `,[L(`> *`,`z-index: 1;`),L(`&::before`,`
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
 `),V(`disabled`,`
 opacity: .45;
 cursor: not-allowed;
 `),V(`collapsed`,[B(`arrow`,`transform: rotate(0);`)]),V(`selected`,[L(`&::before`,`background-color: var(--n-item-color-active);`),B(`arrow`,`color: var(--n-arrow-color-active);`),B(`icon`,`color: var(--n-item-icon-color-active);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-active);
 `,[L(`a`,`color: var(--n-item-text-color-active);`),B(`extra`,`color: var(--n-item-text-color-active);`)])]),V(`child-active`,[F(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active);
 `,[L(`a`,`
 color: var(--n-item-text-color-child-active);
 `),B(`extra`,`
 color: var(--n-item-text-color-child-active);
 `)]),B(`arrow`,`
 color: var(--n-arrow-color-child-active);
 `),B(`icon`,`
 color: var(--n-item-icon-color-child-active);
 `)]),j(`disabled`,[j(`selected, child-active`,[L(`&:focus-within`,He)]),V(`selected`,[Q(null,[B(`arrow`,`color: var(--n-arrow-color-active-hover);`),B(`icon`,`color: var(--n-item-icon-color-active-hover);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-active-hover);
 `,[L(`a`,`color: var(--n-item-text-color-active-hover);`),B(`extra`,`color: var(--n-item-text-color-active-hover);`)])])]),V(`child-active`,[Q(null,[B(`arrow`,`color: var(--n-arrow-color-child-active-hover);`),B(`icon`,`color: var(--n-item-icon-color-child-active-hover);`),F(`menu-item-content-header`,`
 color: var(--n-item-text-color-child-active-hover);
 `,[L(`a`,`color: var(--n-item-text-color-child-active-hover);`),B(`extra`,`color: var(--n-item-text-color-child-active-hover);`)])])]),V(`selected`,[Q(null,[L(`&::before`,`background-color: var(--n-item-color-active-hover);`)])]),Q(null,He)]),B(`icon`,`
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
 `),B(`arrow`,`
 grid-area: arrow;
 font-size: 16px;
 color: var(--n-arrow-color);
 transform: rotate(180deg);
 opacity: 1;
 transition:
 color .3s var(--n-bezier),
 transform 0.2s var(--n-bezier),
 opacity 0.2s var(--n-bezier);
 `),F(`menu-item-content-header`,`
 grid-area: content;
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 opacity: 1;
 white-space: nowrap;
 color: var(--n-item-text-color);
 `,[L(`a`,`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[L(`&::before`,`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),B(`extra`,`
 font-size: .93em;
 color: var(--n-group-text-color);
 transition: color .3s var(--n-bezier);
 `)])]),F(`submenu`,`
 cursor: pointer;
 position: relative;
 margin-top: 6px;
 `,[F(`menu-item-content`,`
 height: var(--n-item-height);
 `),F(`submenu-children`,`
 overflow: hidden;
 padding: 0;
 `,[ge({duration:`.2s`})])]),F(`menu-item-group`,[F(`menu-item-group-title`,`
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
 `)])]),F(`menu-tooltip`,[L(`a`,`
 color: inherit;
 text-decoration: none;
 `)]),F(`menu-divider`,`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function Q(e,t){return[V(`hover`,e,t),L(`&:hover`,e,t)]}var Ge=C({name:`MenuDivider`,setup(){let{mergedClsPrefixRef:t,isHorizontalRef:n}=e(X);return()=>n.value?null:(r(),h(`div`,{key:1,class:E(`${t.value}-menu-divider`)},null,2))}}),Ke=C({name:`ChevronDownFilled`,render(){return(()=>{let e=A(`f3af82a2aab086a5`);return e[0]||=m(`svg`,{viewBox:`0 0 16 16`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},[m(`path`,{d:`M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z`,fill:`currentColor`})],-1)})()}}),qe=[`onClick`],Je=C({name:`MenuOptionContent`,props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(t){let{props:n}=e(X);return{menuProps:n,style:f(()=>{let{paddingLeft:e}=t;return{paddingLeft:e&&`${e}px`}}),iconStyle:f(()=>{let{maxIconSize:e,activeIconSize:n,iconMarginRight:r}=t;return{width:`${e}px`,height:`${e}px`,fontSize:`${n}px`,marginRight:`${r}px`}})}},render(){let{clsPrefix:e,tmNode:t,menuProps:{renderIcon:n,renderLabel:i,renderExtra:a,expandIcon:o}}=this,s=n?n(t.rawNode):W(this.icon);return(()=>{let n=A(`7bb10afc6caf8fa4`);return r(),h(`div`,{onClick:e=>{this.onClick?.(e)},role:`none`,class:E([`${e}-menu-item-content`,{[`${e}-menu-item-content--selected`]:this.selected,[`${e}-menu-item-content--collapsed`]:this.collapsed,[`${e}-menu-item-content--child-active`]:this.childActive,[`${e}-menu-item-content--disabled`]:this.disabled,[`${e}-menu-item-content--hover`]:this.hover}]),style:S(this.style)},[R(()=>s&&(r(),h(`div`,{class:E(`${e}-menu-item-content__icon`),style:S(this.iconStyle),role:`none`},[R(()=>[s])],6))),m(`div`,{class:E(`${e}-menu-item-content-header`),role:`none`},[this.isEllipsisPlaceholder?(r(),h(v,{key:0},[R(()=>this.title)],64)):(r(),h(v,{key:1},[i?(r(),h(v,{key:0},[R(()=>i(t.rawNode))],64)):(r(),h(v,{key:1},[R(()=>W(this.title))],64))],64)),this.extra||a?(r(),h(`span`,{key:2,class:E(`${e}-menu-item-content-header__extra`)},[n[0]||=R(` `,-1),a?(r(),h(v,{key:0},[R(()=>a(t.rawNode))],64)):(r(),h(v,{key:1},[R(()=>W(this.extra))],64))],2)):R(()=>null)],2),this.showArrow?(r(),y(N,{key:0,ariaHidden:!0,class:E(`${e}-menu-item-content__arrow`),clsPrefix:e},{default:()=>o?o(t.rawNode):(r(),y(Ke,{key:1}))},1032,[`class`,`clsPrefix`])):R(()=>null)],14,qe)})()}}),Ye=8;function Xe(t){let n=e(X),{props:r,mergedCollapsedRef:i}=n,a=e(Ve,null),o=e(Z,null),s=f(()=>r.mode===`horizontal`),c=f(()=>s.value?r.dropdownPlacement:`tmNodes`in t?`right-start`:`right`),l=f(()=>Math.max(r.collapsedIconSize??r.iconSize,r.iconSize));return{dropdownPlacement:c,activeIconSize:f(()=>!s.value&&t.root&&i.value?r.collapsedIconSize??r.iconSize:r.iconSize),maxIconSize:l,paddingLeft:f(()=>{if(s.value)return;let{collapsedWidth:e,indent:n,rootIndent:c}=r,{root:u,isGroup:d}=t,f=c===void 0?n:c;return u?i.value?e/2-l.value/2:f:o&&typeof o.paddingLeftRef.value==`number`?i.value?e/2-l.value/2:n/2+o.paddingLeftRef.value:a&&typeof a.paddingLeftRef.value==`number`?(d?n/2:n)+a.paddingLeftRef.value:0}),iconMarginRight:f(()=>{let{collapsedWidth:e,indent:n,rootIndent:a}=r,{value:o}=l,{root:c}=t;return s.value||!c||!i.value?Ye:(a===void 0?n:a)+o+Ye-(e+o)/2}),NMenu:n,NSubmenu:a,NMenuOptionGroup:o}}var Ze={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},Qe={...Ze,tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function},$e=J(Qe),et=C({name:`MenuOption`,props:Qe,setup(e){let t=Xe(e),{NSubmenu:n,NMenu:r,NMenuOptionGroup:i}=t,{props:a,mergedClsPrefixRef:o,mergedCollapsedRef:s}=r,c=n?n.mergedDisabledRef:i?i.mergedDisabledRef:{value:!1},l=f(()=>c.value||e.disabled);function u(t){let{onClick:n}=e;n&&n(t)}function d(t){l.value||(r.doSelect(e.internalKey,e.tmNode.rawNode),u(t))}return{mergedClsPrefix:o,dropdownPlacement:t.dropdownPlacement,paddingLeft:t.paddingLeft,iconMarginRight:t.iconMarginRight,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,mergedTheme:r.mergedThemeRef,menuProps:a,dropdownEnabled:M(()=>e.root&&s.value&&a.mode!==`horizontal`&&!l.value),selected:M(()=>r.mergedValueRef.value===e.internalKey),mergedDisabled:l,handleClick:d}},render(){let{mergedClsPrefix:e,mergedTheme:t,tmNode:n,menuProps:{renderLabel:i,nodeProps:a}}=this,s=a?.(n.rawNode);return r(),h(`div`,o(s,{role:`menuitem`,class:[`${e}-menu-item`,s?.class]}),[(r(),y(me,{theme:t.peers.Tooltip,themeOverrides:t.peerOverrides.Tooltip,trigger:`hover`,placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:[`menu-tooltip`]},{default:()=>i?i(n.rawNode):W(this.title),trigger:()=>(r(),y(Je,{tmNode:n,clsPrefix:e,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick},null,8,[`tmNode`,`clsPrefix`,`paddingLeft`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`selected`,`title`,`extra`,`disabled`,`icon`,`onClick`]))},1032,[`theme`,`themeOverrides`,`placement`,`disabled`]))],16)}}),tt={...Ze,tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}},nt=J(tt),rt=C({name:`MenuOptionGroup`,props:tt,setup(n){let i=Xe(n),{NSubmenu:a}=i,s=f(()=>a?.mergedDisabledRef.value?!0:n.tmNode.disabled);t(Z,{paddingLeftRef:i.paddingLeft,mergedDisabledRef:s});let{mergedClsPrefixRef:c,props:l}=e(X);return function(){let{value:e}=c,t=i.paddingLeft.value,{nodeProps:a}=l,s=a?.(n.tmNode.rawNode);return(()=>{let i=A(`45eca6a63be5028b`);return r(),h(`div`,{class:E(`${e}-menu-item-group`),role:`group`},[m(`div`,o(s,{class:[`${e}-menu-item-group-title`,s?.class],style:[s?.style||``,t===void 0?``:`padding-left: ${t}px;`]}),[R(()=>W(n.title)),n.extra?(r(),h(v,{key:0},[i[0]||=R(` `,-1),R(()=>W(n.extra))],64)):R(()=>null)],16),m(`div`,null,[R(()=>n.tmNodes.map(e=>$(e,l)))])],2)})()}}}),it=[`aria-expanded`,`id`],at=[`aria-expanded`,`id`],ot={...Ze,rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean},st=J(ot),ct=C({name:`Submenu`,props:ot,setup(e){let n=Xe(e),{NMenu:r,NSubmenu:i}=n,{props:a,mergedCollapsedRef:o,mergedThemeRef:s}=r,c=f(()=>{let{disabled:t}=e;return i?.mergedDisabledRef.value||a.disabled?!0:t}),u=l(!1);t(Ve,{paddingLeftRef:n.paddingLeft,mergedDisabledRef:c}),t(Z,null);function d(){let{onClick:t}=e;t&&t()}function p(){c.value||(o.value||r.toggleExpand(e.internalKey),d())}function m(e){u.value=e}return{menuProps:a,mergedTheme:s,doSelect:r.doSelect,inverted:r.invertedRef,isHorizontal:r.isHorizontalRef,mergedClsPrefix:r.mergedClsPrefixRef,maxIconSize:n.maxIconSize,activeIconSize:n.activeIconSize,iconMarginRight:n.iconMarginRight,dropdownPlacement:n.dropdownPlacement,dropdownShow:u,paddingLeft:n.paddingLeft,mergedDisabled:c,mergedValue:r.mergedValueRef,childActive:M(()=>e.virtualChildActive??r.activePathRef.value.includes(e.internalKey)),collapsed:f(()=>a.mode===`horizontal`?!1:o.value?!0:!r.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:f(()=>!c.value&&(a.mode===`horizontal`||o.value)),handlePopoverShowChange:m,handleClick:p}},render(){let{mergedClsPrefix:e,menuProps:{renderIcon:t,renderLabel:n}}=this,i=()=>{let{isHorizontal:e,paddingLeft:t,collapsed:n,mergedDisabled:i,maxIconSize:a,activeIconSize:s,title:c,childActive:l,icon:u,handleClick:d,menuProps:{nodeProps:f},dropdownShow:p,iconMarginRight:m,tmNode:g,mergedClsPrefix:_,isEllipsisPlaceholder:v,extra:b}=this,x=f?.(g.rawNode);return r(),h(`div`,o(x,{class:[`${_}-menu-item`,x?.class],role:`menuitem`}),[(r(),y(Je,{tmNode:g,paddingLeft:t,collapsed:n,disabled:i,iconMarginRight:m,maxIconSize:a,activeIconSize:s,title:c,extra:b,showArrow:!e,childActive:l,clsPrefix:_,icon:u,hover:p,onClick:d,isEllipsisPlaceholder:v},null,8,[`tmNode`,`paddingLeft`,`collapsed`,`disabled`,`iconMarginRight`,`maxIconSize`,`activeIconSize`,`title`,`extra`,`showArrow`,`childActive`,`clsPrefix`,`icon`,`hover`,`onClick`,`isEllipsisPlaceholder`]))],16)},a=()=>(r(),y(ne,null,{default:()=>{let{tmNodes:t,collapsed:n}=this;return n?null:(r(),h(`div`,{key:1,class:E(`${e}-submenu-children`),role:`menu`},[R(()=>t.map(e=>$(e,this.menuProps)))],2))}},1024));return this.root?(r(),y(he,o({key:2,size:`large`,trigger:`hover`},this.menuProps?.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:`14px`,optionIconSizeLarge:`18px`},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:t,renderLabel:n}),{default:()=>(r(),h(`div`,{class:E(`${e}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[R(()=>i()),this.isHorizontal?R(()=>null):(r(),h(v,{key:1},[R(()=>a())],64))],10,it))},1040,[`themeOverrides`,`theme`,`value`,`disabled`,`placement`,`keyField`,`labelField`,`childrenField`,`onUpdateShow`,`options`,`onSelect`,`inverted`,`renderIcon`,`renderLabel`])):(r(),h(`div`,{key:3,class:E(`${e}-submenu`),role:`menu`,"aria-expanded":!this.collapsed,id:this.domId},[R(()=>i()),R(()=>a())],10,at))}});function lt(e){return e.type===`divider`||e.type===`render`}function ut(e){return e.type===`divider`}function $(e,t){let{rawNode:n}=e,{show:i}=n;if(i===!1)return null;if(lt(n))return ut(n)?(r(),y(Ge,o({key:e.key},n.props),null,16)):null;let{labelField:s}=t,{key:c,level:l,isGroup:u}=e,d={...n,title:n.title||n[s],extra:n.titleExtra||n.extra,key:c,internalKey:c,level:l,root:l===0,isGroup:u};return e.children?e.isGroup?a(rt,H(d,nt,{tmNode:e,tmNodes:e.children,key:c})):a(ct,H(d,st,{key:c,rawNodes:n[t.childrenField],tmNodes:e.children,tmNode:e})):a(et,H(d,$e,{key:c,tmNode:e}))}var dt={...z.props,options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:`label`},keyField:{type:String,default:`key`},childrenField:{type:String,default:`children`},disabledField:{type:String,default:`disabled`},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:`vertical`},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:`bottom`},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array},ft=C({name:`Menu`,inheritAttrs:!1,props:dt,setup(n){let{mergedClsPrefixRef:i,inlineThemeDisabled:a}=te(n),o=z(`Menu`,`-menu`,We,ye,n,i),c=e(Te,null),u=f(()=>{let{collapsed:e}=n;if(e!==void 0)return e;if(c){let{collapseModeRef:e,collapsedRef:t}=c;if(e.value===`width`)return t.value??!1}return!1}),d=f(()=>{let{keyField:e,childrenField:t,disabledField:r}=n;return G(n.items||n.options,{getIgnored(e){return lt(e)},getChildren(e){return e[t]},getDisabled(e){return e[r]},getKey(t){return t[e]??t.name}})}),m=f(()=>new Set(d.value.treeNodes.map(e=>e.key))),{watchProps:h}=n,g=l(null);h?.includes(`defaultValue`)?s(()=>{g.value=n.defaultValue}):g.value=n.defaultValue;let _=p(n,`value`),v=U(_,g),b=l([]),x=()=>{b.value=n.defaultExpandAll?d.value.getNonLeafKeys():n.defaultExpandedNames||n.defaultExpandedKeys||d.value.getPath(v.value,{includeSelf:!1}).keyPath};h?.includes(`defaultExpandedKeys`)?s(x):x();let S=oe(n,[`expandedNames`,`expandedKeys`]),C=U(S,b),w=f(()=>d.value.treeNodes),T=f(()=>d.value.getPath(v.value).keyPath);t(X,{props:n,mergedCollapsedRef:u,mergedThemeRef:o,mergedValueRef:v,mergedExpandedKeysRef:C,activePathRef:T,mergedClsPrefixRef:i,isHorizontalRef:f(()=>n.mode===`horizontal`),invertedRef:p(n,`inverted`),doSelect:E,toggleExpand:O});function E(e,t){let{"onUpdate:value":r,onUpdateValue:i,onSelect:a}=n;i&&q(i,e,t),r&&q(r,e,t),a&&q(a,e,t),g.value=e}function D(e){let{"onUpdate:expandedKeys":t,onUpdateExpandedKeys:r,onExpandedNamesChange:i,onOpenNamesChange:a}=n;t&&q(t,e),r&&q(r,e),i&&q(i,e),a&&q(a,e),b.value=e}function O(e){let t=Array.from(C.value),r=t.findIndex(t=>t===e);if(~r)t.splice(r,1);else{if(n.accordion&&m.value.has(e)){let e=t.findIndex(e=>m.value.has(e));e>-1&&t.splice(e,1)}t.push(e)}D(t)}let ee=e=>{let t=d.value.getPath(e??v.value,{includeSelf:!1}).keyPath;if(!t.length)return;let r=Array.from(C.value),i=new Set([...r,...t]);n.accordion&&m.value.forEach(e=>{i.has(e)&&!t.includes(e)&&i.delete(e)}),D(Array.from(i))},k=f(()=>{let{inverted:e}=n,{common:{cubicBezierEaseInOut:t},self:r}=o.value,{borderRadius:i,borderColorHorizontal:a,fontSize:s,itemHeight:c,dividerColor:l}=r,u={"--n-divider-color":l,"--n-bezier":t,"--n-font-size":s,"--n-border-color-horizontal":a,"--n-border-radius":i,"--n-item-height":c};return e?(u[`--n-group-text-color`]=r.groupTextColorInverted,u[`--n-color`]=r.colorInverted,u[`--n-item-text-color`]=r.itemTextColorInverted,u[`--n-item-text-color-hover`]=r.itemTextColorHoverInverted,u[`--n-item-text-color-active`]=r.itemTextColorActiveInverted,u[`--n-item-text-color-child-active`]=r.itemTextColorChildActiveInverted,u[`--n-item-text-color-child-active-hover`]=r.itemTextColorChildActiveInverted,u[`--n-item-text-color-active-hover`]=r.itemTextColorActiveHoverInverted,u[`--n-item-icon-color`]=r.itemIconColorInverted,u[`--n-item-icon-color-hover`]=r.itemIconColorHoverInverted,u[`--n-item-icon-color-active`]=r.itemIconColorActiveInverted,u[`--n-item-icon-color-active-hover`]=r.itemIconColorActiveHoverInverted,u[`--n-item-icon-color-child-active`]=r.itemIconColorChildActiveInverted,u[`--n-item-icon-color-child-active-hover`]=r.itemIconColorChildActiveHoverInverted,u[`--n-item-icon-color-collapsed`]=r.itemIconColorCollapsedInverted,u[`--n-item-text-color-horizontal`]=r.itemTextColorHorizontalInverted,u[`--n-item-text-color-hover-horizontal`]=r.itemTextColorHoverHorizontalInverted,u[`--n-item-text-color-active-horizontal`]=r.itemTextColorActiveHorizontalInverted,u[`--n-item-text-color-child-active-horizontal`]=r.itemTextColorChildActiveHorizontalInverted,u[`--n-item-text-color-child-active-hover-horizontal`]=r.itemTextColorChildActiveHoverHorizontalInverted,u[`--n-item-text-color-active-hover-horizontal`]=r.itemTextColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-horizontal`]=r.itemIconColorHorizontalInverted,u[`--n-item-icon-color-hover-horizontal`]=r.itemIconColorHoverHorizontalInverted,u[`--n-item-icon-color-active-horizontal`]=r.itemIconColorActiveHorizontalInverted,u[`--n-item-icon-color-active-hover-horizontal`]=r.itemIconColorActiveHoverHorizontalInverted,u[`--n-item-icon-color-child-active-horizontal`]=r.itemIconColorChildActiveHorizontalInverted,u[`--n-item-icon-color-child-active-hover-horizontal`]=r.itemIconColorChildActiveHoverHorizontalInverted,u[`--n-arrow-color`]=r.arrowColorInverted,u[`--n-arrow-color-hover`]=r.arrowColorHoverInverted,u[`--n-arrow-color-active`]=r.arrowColorActiveInverted,u[`--n-arrow-color-active-hover`]=r.arrowColorActiveHoverInverted,u[`--n-arrow-color-child-active`]=r.arrowColorChildActiveInverted,u[`--n-arrow-color-child-active-hover`]=r.arrowColorChildActiveHoverInverted,u[`--n-item-color-hover`]=r.itemColorHoverInverted,u[`--n-item-color-active`]=r.itemColorActiveInverted,u[`--n-item-color-active-hover`]=r.itemColorActiveHoverInverted,u[`--n-item-color-active-collapsed`]=r.itemColorActiveCollapsedInverted):(u[`--n-group-text-color`]=r.groupTextColor,u[`--n-color`]=r.color,u[`--n-item-text-color`]=r.itemTextColor,u[`--n-item-text-color-hover`]=r.itemTextColorHover,u[`--n-item-text-color-active`]=r.itemTextColorActive,u[`--n-item-text-color-child-active`]=r.itemTextColorChildActive,u[`--n-item-text-color-child-active-hover`]=r.itemTextColorChildActiveHover,u[`--n-item-text-color-active-hover`]=r.itemTextColorActiveHover,u[`--n-item-icon-color`]=r.itemIconColor,u[`--n-item-icon-color-hover`]=r.itemIconColorHover,u[`--n-item-icon-color-active`]=r.itemIconColorActive,u[`--n-item-icon-color-active-hover`]=r.itemIconColorActiveHover,u[`--n-item-icon-color-child-active`]=r.itemIconColorChildActive,u[`--n-item-icon-color-child-active-hover`]=r.itemIconColorChildActiveHover,u[`--n-item-icon-color-collapsed`]=r.itemIconColorCollapsed,u[`--n-item-text-color-horizontal`]=r.itemTextColorHorizontal,u[`--n-item-text-color-hover-horizontal`]=r.itemTextColorHoverHorizontal,u[`--n-item-text-color-active-horizontal`]=r.itemTextColorActiveHorizontal,u[`--n-item-text-color-child-active-horizontal`]=r.itemTextColorChildActiveHorizontal,u[`--n-item-text-color-child-active-hover-horizontal`]=r.itemTextColorChildActiveHoverHorizontal,u[`--n-item-text-color-active-hover-horizontal`]=r.itemTextColorActiveHoverHorizontal,u[`--n-item-icon-color-horizontal`]=r.itemIconColorHorizontal,u[`--n-item-icon-color-hover-horizontal`]=r.itemIconColorHoverHorizontal,u[`--n-item-icon-color-active-horizontal`]=r.itemIconColorActiveHorizontal,u[`--n-item-icon-color-active-hover-horizontal`]=r.itemIconColorActiveHoverHorizontal,u[`--n-item-icon-color-child-active-horizontal`]=r.itemIconColorChildActiveHorizontal,u[`--n-item-icon-color-child-active-hover-horizontal`]=r.itemIconColorChildActiveHoverHorizontal,u[`--n-arrow-color`]=r.arrowColor,u[`--n-arrow-color-hover`]=r.arrowColorHover,u[`--n-arrow-color-active`]=r.arrowColorActive,u[`--n-arrow-color-active-hover`]=r.arrowColorActiveHover,u[`--n-arrow-color-child-active`]=r.arrowColorChildActive,u[`--n-arrow-color-child-active-hover`]=r.arrowColorChildActiveHover,u[`--n-item-color-hover`]=r.itemColorHover,u[`--n-item-color-active`]=r.itemColorActive,u[`--n-item-color-active-hover`]=r.itemColorActiveHover,u[`--n-item-color-active-collapsed`]=r.itemColorActiveCollapsed),u}),A=a?I(`menu`,f(()=>n.inverted?`a`:`b`),k,n):void 0,j=ae(),M=l(null),N=l(null),P=!0,F=()=>{P?P=!1:M.value?.sync({showAllItemsBeforeCalculate:!0})};function L(){return document.getElementById(j)}let R=l(-1);function ne(e){R.value=n.options.length-e}function re(e){e||(R.value=-1)}let ie=f(()=>{let e=R.value;return{children:e===-1?[]:n.options.slice(e)}}),B=f(()=>{let{childrenField:e,disabledField:t,keyField:r}=n;return G([ie.value],{getIgnored(e){return lt(e)},getChildren(t){return t[e]},getDisabled(e){return e[t]},getKey(e){return e[r]??e.name}})}),V=f(()=>G([{}]).treeNodes[0]);function H(){if(R.value===-1)return r(),y(ct,{root:!0,level:0,key:`__ellpisisGroupPlaceholder__`,internalKey:`__ellpisisGroupPlaceholder__`,title:`···`,tmNode:V.value,domId:j,isEllipsisPlaceholder:!0},null,8,[`tmNode`,`domId`]);let e=B.value.treeNodes[0],t=T.value,n=!!e.children?.some(e=>t.includes(e.key));return r(),y(ct,{level:0,root:!0,key:`__ellpisisGroup__`,internalKey:`__ellpisisGroup__`,title:`···`,virtualChildActive:n,tmNode:e,domId:j,rawNodes:e.rawNode.children||[],tmNodes:e.children||[],isEllipsisPlaceholder:!0},null,8,[`virtualChildActive`,`tmNode`,`domId`,`rawNodes`,`tmNodes`])}return{mergedClsPrefix:i,controlledExpandedKeys:S,uncontrolledExpanededKeys:b,mergedExpandedKeys:C,uncontrolledValue:g,mergedValue:v,activePath:T,tmNodes:w,mergedTheme:o,mergedCollapsed:u,cssVars:a?void 0:k,themeClass:A?.themeClass,overflowRef:M,counterRef:N,updateCounter:()=>{},onResize:F,onUpdateOverflow:re,onUpdateCount:ne,renderCounter:H,getCounter:L,onRender:A?.onRender,showOption:ee,deriveResponsiveState:F}},render(){let{mergedClsPrefix:e,mode:t,themeClass:n,onRender:i}=this;i?.();let s=()=>this.tmNodes.map(e=>$(e,this.$props)),c=t===`horizontal`&&this.responsive,l=()=>a(`div`,o(this.$attrs,{role:t===`horizontal`?`menubar`:`menu`,class:[`${e}-menu`,n,`${e}-menu--${t}`,c&&`${e}-menu--responsive`,this.mergedCollapsed&&`${e}-menu--collapsed`],style:this.cssVars}),c?(r(),y(se,{key:2,ref:`overflowRef`,onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:this.renderCounter},1032,[`onUpdateOverflow`,`getCounter`,`onUpdateCount`,`updateCounter`])):s());return c?(r(),y(ue,{key:3,onResize:this.onResize},{default:l},1032,[`onResize`])):l()}}),pt={class:`project-switcher`},mt={key:1,class:`none-hint`},ht={key:3,class:`quick-create`},gt=[`placeholder`],_t=Se(C({__name:`ProjectSwitcher`,setup(e){let{t}=_(),{projects:i,currentId:a,select:o,create:s}=xe(),d=l(!1),f=l(``);async function p(){f.value.trim()&&(await s(f.value.trim(),``),f.value=``,d.value=!1)}return(e,s)=>(r(),h(`div`,pt,[b(i).length>0?(r(),y(b(ce),{key:0,value:b(a),options:b(i).map(e=>({label:e.name,value:e.id})),size:`small`,placeholder:b(t)(`nav.projectNone`),"onUpdate:value":s[0]||=e=>b(o)(e)},null,8,[`value`,`options`,`placeholder`])):(r(),h(`div`,mt,w(b(t)(`nav.projectNone`)),1)),d.value?(r(),h(`div`,ht,[n(m(`input`,{"onUpdate:modelValue":s[2]||=e=>f.value=e,class:`n-input`,placeholder:b(t)(`project.name`),onKeyup:[T(p,[`enter`]),s[3]||=T(e=>d.value=!1,[`esc`])]},null,40,gt),[[ee,f.value]]),x(b(Y),{size:`tiny`,type:`primary`,onClick:p},{default:c(()=>[...s[5]||=[u(`OK`,-1)]]),_:1})])):(r(),y(b(Y),{key:2,quaternary:``,size:`tiny`,class:`add-btn`,title:b(t)(`project.new`),onClick:s[1]||=e=>d.value=!0},{default:c(()=>[...s[4]||=[u(`+`,-1)]]),_:1},8,[`title`]))]))}}),[[`__scopeId`,`data-v-140fd69a`]]),vt={class:`brand`},yt={class:`brand-text`},bt={class:`brand-sub`},xt={class:`sider-footer`},St={class:`content-scroll`},Ct=Se(C({__name:`LayoutView`,setup(e){let{t,locale:n}=_(),a=ve(),{load:o}=xe();g(o);let s=f(()=>[{label:t(`nav.dashboard`),key:`/`},{label:t(`nav.build`),key:`/build/items`},{label:t(`nav.deploy`),key:`/deploy/devices`},{label:t(`nav.test`),key:`/test`},{label:t(`nav.release`),key:`/release`},{label:t(`nav.settings`),key:`/settings`}]);function l(e){a.push(e)}function p(){n.value=n.value===`zh`?`en`:`zh`}function h(){de(),a.push(`/login`)}return(e,n)=>{let a=i(`router-view`);return r(),y(b(je),{position:`absolute`,class:d([`root`,{dark:b(be).isDark}]),"has-sider":``},{default:c(()=>[x(b(Be),{bordered:``,class:`sider`},{default:c(()=>[m(`div`,vt,[n[1]||=m(`div`,{class:`brand-mark`},`EF`,-1),m(`div`,yt,[n[0]||=m(`div`,{class:`brand-name`},`EmbedFlow`,-1),m(`div`,bt,w(b(t)(`project.title`)),1)])]),x(_t),x(b(ft),{options:s.value,"onUpdate:value":l,class:`nav-menu`},null,8,[`options`]),m(`div`,xt,[x(b(Y),{quaternary:``,size:`small`,onClick:b(_e)},{default:c(()=>[u(w(b(be).isDark?`Light`:`Dark`),1)]),_:1},8,[`onClick`]),x(b(Y),{quaternary:``,size:`small`,onClick:p},{default:c(()=>[u(w(b(t)(`nav.language`)),1)]),_:1}),x(b(Y),{quaternary:``,size:`small`,onClick:h},{default:c(()=>[u(w(b(t)(`nav.logout`)),1)]),_:1})])]),_:1}),x(b(Me),{class:`content-outer`},{default:c(()=>[m(`div`,St,[x(a)])]),_:1})]),_:1},8,[`class`])}}}),[[`__scopeId`,`data-v-fbf39fff`]]);export{Ct as default};