import{C as e,F as t,G as n,H as r,P as i,S as a,T as o,U as s,X as c,d as l,et as u,f as d,h as f,k as p,o as m,p as h,vt as g,y as _}from"./vue-i18n-CU1juWHN.js";import{E as v,Et as y,M as b,O as x,P as S,R as C,gt as w,h as T,ht as E,k as D,kt as O,l as k,m as A,s as j,u as M,vt as N,yt as P}from"./light-uQ0rL05w.js";import{E as F,S as I,T as L,c as R,d as ee,l as z,t as B,u as V,v as H,x as U}from"./event-D8DEXtwy.js";import{t as W}from"./format-length-IO90KiIs.js";import{N as G,T as K,x as q}from"./http-gZvvjYrt.js";import{i as J,n as Y}from"./Modal-4eTFlEpF.js";import{l as X}from"./index-C20PsuhH.js";var Z=[`onMouseenter`,`onMouseleave`,`onMousedown`],Q={key:1,role:`none`},te=_({name:`NDrawerContent`,inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(n){let i=c(!!n.show),a=c(null),o=e(F),u=0,d=``,f=null,m=c(!1),h=c(!1),g=l(()=>n.placement===`top`||n.placement===`bottom`),{mergedClsPrefixRef:_,mergedRtlRef:v}=S(n),y=k(`Drawer`,v,_),b=M,x=e=>{h.value=!0,u=g.value?e.clientY:e.clientX,d=document.body.style.cursor,document.body.style.cursor=g.value?`ns-resize`:`ew-resize`,document.body.addEventListener(`mousemove`,j),document.body.addEventListener(`mouseleave`,b),document.body.addEventListener(`mouseup`,M)},C=()=>{f!==null&&(window.clearTimeout(f),f=null),h.value?m.value=!0:f=window.setTimeout(()=>{m.value=!0},300)},w=()=>{f!==null&&(window.clearTimeout(f),f=null),m.value=!1},{doUpdateHeight:T,doUpdateWidth:E}=o,D=e=>{let{maxWidth:t}=n;if(t&&e>t)return t;let{minWidth:r}=n;return r&&e<r?r:e},A=e=>{let{maxHeight:t}=n;if(t&&e>t)return t;let{minHeight:r}=n;return r&&e<r?r:e};function j(e){if(h.value){if(g.value){let t=a.value?.offsetHeight||0,r=u-e.clientY;t+=n.placement===`bottom`?r:-r,t=A(t),T(t),u=e.clientY}else{let t=a.value?.offsetWidth||0,r=u-e.clientX;t+=n.placement===`right`?r:-r,t=D(t),E(t),u=e.clientX}}}function M(){h.value&&(u=0,h.value=!1,document.body.style.cursor=d,document.body.removeEventListener(`mousemove`,j),document.body.removeEventListener(`mouseup`,M),document.body.removeEventListener(`mouseleave`,b))}s(()=>{n.show&&(i.value=!0)}),r(()=>n.show,e=>{e||M()}),p(()=>{M()});let N=l(()=>{let{show:e}=n,t=[[O,e]];return n.showMask||t.push([ee,n.onClickoutside,void 0,{capture:!0}]),t});function P(){i.value=!1,n.onAfterLeave?.()}return J(l(()=>n.blockScroll&&i.value)),t(L,a),t(U,null),t(I,null),{bodyRef:a,rtlEnabled:y,mergedClsPrefix:o.mergedClsPrefixRef,isMounted:o.isMountedRef,mergedTheme:o.mergedThemeRef,displayed:i,transitionName:l(()=>({right:`slide-in-from-right-transition`,left:`slide-in-from-left-transition`,top:`slide-in-from-top-transition`,bottom:`slide-in-from-bottom-transition`})[n.placement]),handleAfterLeave:P,bodyDirectives:N,handleMousedownResizeTrigger:x,handleMouseenterResizeTrigger:C,handleMouseleaveResizeTrigger:w,isDragging:h,isHoverOnResizeTrigger:m}},render(){let{$slots:e,mergedClsPrefix:t}=this;return this.displayDirective===`show`||this.displayed||this.show?n((i(),f(`div`,Q,[(i(),h(R,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>(i(),h(y,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>n(a(`div`,o(this.$attrs,{role:`dialog`,ref:`bodyRef`,"aria-modal":`true`,class:[`${t}-drawer`,this.rtlEnabled&&`${t}-drawer--rtl`,`${t}-drawer--${this.placement}-placement`,this.isDragging&&`${t}-drawer--unselectable`,this.nativeScrollbar&&`${t}-drawer--native-scrollbar`]}),[this.resizable?(i(),f(`div`,{key:2,class:v([`${t}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${t}-drawer__resize-trigger--hover`]),onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger},null,42,Z)):null,this.nativeScrollbar?(i(),f(`div`,{key:3,class:v([`${t}-drawer-content-wrapper`,this.contentClass]),style:g(this.contentStyle),role:`none`},[D(()=>e.default?.())],6)):(i(),h(q,o({key:4},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${t}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),x(e),1040,[`contentStyle`,`contentClass`,`theme`,`themeOverrides`]))]),this.bodyDirectives)},1032,[`name`,`appear`,`onAfterEnter`,`onAfterLeave`]))},1032,[`disabled`,`active`,`autoFocus`,`onEsc`]))])),[[O,this.displayDirective===`if`||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:ne,cubicBezierEaseOut:re}=b;function ie({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-bottom`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${ne}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${re}`}),E(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateY(100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateY(100%)`})]}var{cubicBezierEaseIn:ae,cubicBezierEaseOut:oe}=b;function se({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-left`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${ae}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${oe}`}),E(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateX(-100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateX(-100%)`})]}var{cubicBezierEaseIn:ce,cubicBezierEaseOut:le}=b;function ue({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-right`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${ce}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${le}`}),E(`&.${n}-transition-enter-to`,{transform:`translateX(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateX(100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateX(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateX(100%)`})]}var{cubicBezierEaseIn:de,cubicBezierEaseOut:fe}=b;function pe({duration:e=`0.3s`,leaveDuration:t=`0.2s`,name:n=`slide-in-from-top`}={}){return[E(`&.${n}-transition-leave-active`,{transition:`transform ${t} ${de}`}),E(`&.${n}-transition-enter-active`,{transition:`transform ${e} ${fe}`}),E(`&.${n}-transition-enter-to`,{transform:`translateY(0)`}),E(`&.${n}-transition-enter-from`,{transform:`translateY(-100%)`}),E(`&.${n}-transition-leave-from`,{transform:`translateY(0)`}),E(`&.${n}-transition-leave-to`,{transform:`translateY(-100%)`})]}var $=E([w(`drawer`,`
 word-break: break-word;
 line-height: var(--n-line-height);
 position: absolute;
 pointer-events: all;
 box-shadow: var(--n-box-shadow);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 background-color: var(--n-color);
 color: var(--n-text-color);
 box-sizing: border-box;
 `,[ue(),se(),pe(),ie(),P(`unselectable`,`
 user-select: none; 
 -webkit-user-select: none;
 `),P(`native-scrollbar`,[w(`drawer-content-wrapper`,`
 overflow: auto;
 height: 100%;
 `)]),N(`resize-trigger`,`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[P(`hover`,`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),w(`drawer-content-wrapper`,`
 box-sizing: border-box;
 `),w(`drawer-content`,`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[P(`native-scrollbar`,[w(`drawer-body-content-wrapper`,`
 height: 100%;
 overflow: auto;
 `)]),w(`drawer-body`,`
 flex: 1 0 0;
 overflow: hidden;
 `),w(`drawer-body-content-wrapper`,`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),w(`drawer-header`,`
 font-weight: var(--n-title-font-weight);
 line-height: 1;
 font-size: var(--n-title-font-size);
 color: var(--n-title-text-color);
 padding: var(--n-header-padding);
 transition: border .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-divider-color);
 border-bottom: var(--n-header-border-bottom);
 display: flex;
 justify-content: space-between;
 align-items: center;
 `,[N(`main`,`
 flex: 1;
 `),N(`close`,`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),w(`drawer-footer`,`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),P(`right-placement`,`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),P(`left-placement`,`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),P(`top-placement`,`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),P(`bottom-placement`,`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[N(`resize-trigger`,`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),E(`body`,[E(`>`,[w(`drawer-container`,`
 position: fixed;
 `)])]),w(`drawer-container`,`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[E(`> *`,`
 pointer-events: all;
 `)]),w(`drawer-mask`,`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[P(`invisible`,`
 background-color: rgba(0, 0, 0, 0)
 `),K({enterDuration:`0.2s`,leaveDuration:`0.2s`,enterCubicBezier:`var(--n-bezier-in)`,leaveCubicBezier:`var(--n-bezier-out)`})])]),me=[`onClick`],he={...A.props,show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:`right`},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:`if`},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function},ge=_({name:`Drawer`,inheritAttrs:!1,props:he,setup(e){let{mergedClsPrefixRef:n,namespaceRef:r,inlineThemeDisabled:i}=S(e),a=M(),o=A(`Drawer`,`-drawer`,$,X,e,n),s=c(e.defaultWidth),d=c(e.defaultHeight),f=H(u(e,`width`),s),p=H(u(e,`height`),d),m=l(()=>{let{placement:t}=e;return t===`top`||t===`bottom`?``:W(f.value)}),h=l(()=>{let{placement:t}=e;return t===`left`||t===`right`?``:W(p.value)}),g=t=>{let{onUpdateWidth:n,"onUpdate:width":r}=e;n&&G(n,t),r&&G(r,t),s.value=t},_=t=>{let{onUpdateHeight:n,"onUpdate:width":r}=e;n&&G(n,t),r&&G(r,t),d.value=t},v=l(()=>[{width:m.value,height:h.value},e.drawerStyle||``]);function y(t){let{onMaskClick:n,maskClosable:r}=e;r&&w(!1),n&&n(t)}function b(e){y(e)}let x=Y();function C(t){e.onEsc?.(),e.show&&e.closeOnEsc&&B(t)&&(x.value||w(!1))}function w(t){let{onHide:n,onUpdateShow:r,"onUpdate:show":i}=e;r&&G(r,t),i&&G(i,t),n&&!t&&G(n,t)}t(F,{isMountedRef:a,mergedThemeRef:o,mergedClsPrefixRef:n,doUpdateShow:w,doUpdateHeight:_,doUpdateWidth:g});let E=l(()=>{let{common:{cubicBezierEaseInOut:e,cubicBezierEaseIn:t,cubicBezierEaseOut:n},self:{color:r,textColor:i,boxShadow:a,lineHeight:s,headerPadding:c,footerPadding:l,borderRadius:u,bodyPadding:d,titleFontSize:f,titleTextColor:p,titleFontWeight:m,headerBorderBottom:h,footerBorderTop:g,closeIconColor:_,closeIconColorHover:v,closeIconColorPressed:y,closeColorHover:b,closeColorPressed:x,closeIconSize:S,closeSize:C,closeBorderRadius:w,resizableTriggerColorHover:T}}=o.value;return{"--n-line-height":s,"--n-color":r,"--n-border-radius":u,"--n-text-color":i,"--n-box-shadow":a,"--n-bezier":e,"--n-bezier-out":n,"--n-bezier-in":t,"--n-header-padding":c,"--n-body-padding":d,"--n-footer-padding":l,"--n-title-text-color":p,"--n-title-font-size":f,"--n-title-font-weight":m,"--n-header-border-bottom":h,"--n-footer-border-top":g,"--n-close-icon-color":_,"--n-close-icon-color-hover":v,"--n-close-icon-color-pressed":y,"--n-close-size":C,"--n-close-color-hover":b,"--n-close-color-pressed":x,"--n-close-icon-size":S,"--n-close-border-radius":w,"--n-resize-trigger-color-hover":T}}),D=i?T(`drawer`,void 0,E,e):void 0;return{mergedClsPrefix:n,namespace:r,mergedBodyStyle:v,handleOutsideClick:b,handleMaskClick:y,handleEsc:C,mergedTheme:o,cssVars:i?void 0:E,themeClass:D?.themeClass,onRender:D?.onRender,isMounted:a}},render(){let{mergedClsPrefix:e}=this;return i(),h(z,{to:this.to,show:this.show},{default:()=>(this.onRender?.(),n((i(),f(`div`,{class:v([`${e}-drawer-container`,this.namespace,this.themeClass]),style:g(this.cssVars),role:`none`},[this.showMask?(i(),h(y,{key:0,name:`fade-in-transition`,appear:this.isMounted},{default:()=>this.show?(i(),f(`div`,{key:1,"aria-hidden":!0,class:v([`${e}-drawer-mask`,this.showMask===`transparent`&&`${e}-drawer-mask--invisible`]),onClick:this.handleMaskClick},null,10,me)):null},1032,[`appear`])):D(()=>null),(i(),h(te,o(this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),x(this.$slots),1040,[`class`,`style`,`blockScroll`,`contentStyle`,`contentClass`,`placement`,`scrollbarProps`,`show`,`displayDirective`,`nativeScrollbar`,`onAfterEnter`,`onAfterLeave`,`trapFocus`,`autoFocus`,`resizable`,`maxHeight`,`minHeight`,`maxWidth`,`minWidth`,`showMask`,`onEsc`,`onClickoutside`]))],6)),[[V,{zIndex:this.zIndex,enabled:this.show}]]))},1032,[`to`,`show`])}}),_e=_({name:`DrawerContent`,props:{title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},slots:Object,setup(){let t=e(F,null);t||C(`drawer-content`,"`n-drawer-content` must be placed inside `n-drawer`.");let{doUpdateShow:n}=t;function r(){n(!1)}return{handleCloseClick:r,mergedTheme:t.mergedThemeRef,mergedClsPrefix:t.mergedClsPrefixRef}},render(){let{title:e,mergedClsPrefix:t,nativeScrollbar:n,mergedTheme:r,bodyClass:a,bodyStyle:s,bodyContentClass:c,bodyContentStyle:l,headerClass:u,headerStyle:p,footerClass:_,footerStyle:y,scrollbarProps:b,closable:S,$slots:C}=this;return i(),f(`div`,{role:`none`,class:v([`${t}-drawer-content`,n&&`${t}-drawer-content--native-scrollbar`])},[C.header||e||S?(i(),f(`div`,{key:0,class:v([`${t}-drawer-header`,u]),style:g(p),role:`none`},[d(`div`,{class:v(`${t}-drawer-header__main`),role:`heading`,"aria-level":`1`},[C.header===void 0?(i(),f(m,{key:1},[D(()=>e)],64)):(i(),f(m,{key:0},[D(()=>C.header())],64))],2),D(()=>S&&(i(),h(j,{onClick:this.handleCloseClick,clsPrefix:t,class:v(`${t}-drawer-header__close`),absolute:!0},null,8,[`onClick`,`clsPrefix`,`class`])))],6)):D(()=>null),n?(i(),f(`div`,{key:2,class:v([`${t}-drawer-body`,a]),style:g(s),role:`none`},[d(`div`,{class:v([`${t}-drawer-body-content-wrapper`,c]),style:g(l),role:`none`},[D(()=>C.default?.())],6)],6)):(i(),h(q,o({key:3,themeOverrides:r.peerOverrides.Scrollbar,theme:r.peers.Scrollbar},b,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,c],contentStyle:l}),x(C),1040,[`themeOverrides`,`theme`,`class`,`contentClass`,`contentStyle`])),C.footer?(i(),f(`div`,{key:4,class:v([`${t}-drawer-footer`,_]),style:g(y),role:`none`},[D(()=>C.footer())],6)):D(()=>null)],2)}});export{ge as n,_e as t};