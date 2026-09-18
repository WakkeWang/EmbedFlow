import{Bt as e,E as t,Et as n,Ft as r,It as i,Jt as a,Kt as o,L as s,Lt as c,Mt as l,O as u,P as d,Qt as f,Rt as p,St as m,Ut as h,Yt as g,Zt as _,bn as v,f as y,gt as b,h as x,hn as S,ht as C,k as w,kt as T,l as E,m as D,mn as O,mt as k,nn as A,rn as j,sn as M,tn as ee,u as te,un as N,vn as P,vt as F,xn as I,xt as ne,yt as L,z as re}from"./vue-i18n-BxOqVJSv.js";import{t as ie}from"./event-Ef8K3FrB.js";import{A as R,D as z,F as ae,I as B,N as oe,O as V,T as H,b as se,g as ce,k as U,l as W}from"./_plugin-vue_export-helper-3sd4LcYN.js";import{_ as G,a as le,c as K,g as ue,h as de,i as fe,m as pe,n as me,p as he,r as ge,t as _e,u as ve,v as ye,y as be}from"./fade-in-scale-up.cssr-DpVYUTqJ.js";import{r as xe,t as Se}from"./use-is-composing-YehLS6Fp.js";import{n as Ce,r as we,t as Te}from"./Card-dMVqAluv.js";import{A as Ee,E as q,M as De,N as Oe,d as ke,f as Ae,j as je,k as Me}from"./index-0G3VD5Ph.js";var J=S(null);function Ne(e){if(e.clientX>0||e.clientY>0)J.value={x:e.clientX,y:e.clientY};else{let{target:t}=e;if(t instanceof Element){let{left:e,top:n,width:r,height:i}=t.getBoundingClientRect();J.value=e>0||n>0?{x:e+r/2,y:n+i/2}:{x:0,y:0}}else J.value=null}}var Y=0,Pe=!0;function Fe(){if(!pe)return O(S(null));Y===0&&R(`click`,document,Ne,!0);let e=()=>{Y+=1};return(Pe&&=he())?(_(e),f(()=>{--Y,Y===0&&U(`click`,document,Ne,!0)})):e(),O(J)}var Ie=S(void 0),X=0;function Le(){Ie.value=Date.now()}var Re=!0;function ze(e){if(!pe)return O(S(!1));let t=S(!1),n=null;function r(){n!==null&&window.clearTimeout(n)}function i(){r(),t.value=!0,n=window.setTimeout(()=>{t.value=!1},e)}X===0&&R(`click`,window,Le,!0);let a=()=>{X+=1,R(`click`,window,i,!0)};return(Re&&=he())?(_(a),f(()=>{--X,X===0&&U(`click`,window,Le,!0),U(`click`,window,i,!0),r()})):a(),O(t)}var Be=s(`n-dialog-provider`);s(`n-dialog-api`),s(`n-dialog-reactive-list`);var Z={icon:Function,type:{type:String,default:`default`},title:[String,Function],closable:{type:Boolean,default:!0},negativeText:String,positiveText:String,positiveButtonProps:Object,negativeButtonProps:Object,content:[String,Function],action:Function,showIcon:{type:Boolean,default:!0},loading:Boolean,bordered:Boolean,iconPlacement:String,titleClass:[String,Array],titleStyle:[String,Object],contentClass:[String,Array],contentStyle:[String,Object],actionClass:[String,Array],actionStyle:[String,Object],onPositiveClick:Function,onNegativeClick:Function,onClose:Function,closeFocusable:Boolean},Ve=B(Z),He=C([b(`dialog`,`
 --n-icon-margin: var(--n-icon-margin-top) var(--n-icon-margin-right) var(--n-icon-margin-bottom) var(--n-icon-margin-left);
 word-break: break-word;
 line-height: var(--n-line-height);
 position: relative;
 background: var(--n-color);
 color: var(--n-text-color);
 box-sizing: border-box;
 margin: auto;
 border-radius: var(--n-border-radius);
 padding: var(--n-padding);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `,[F(`icon`,`
 color: var(--n-icon-color);
 `),L(`bordered`,`
 border: var(--n-border);
 `),L(`icon-top`,[F(`close`,`
 margin: var(--n-close-margin);
 `),F(`icon`,`
 margin: var(--n-icon-margin);
 `),F(`content`,`
 text-align: center;
 `),F(`title`,`
 justify-content: center;
 `),F(`action`,`
 justify-content: center;
 `)]),L(`icon-left`,[F(`icon`,`
 margin: var(--n-icon-margin);
 `),L(`closable`,[F(`title`,`
 padding-right: calc(var(--n-close-size) + 6px);
 `)])]),F(`close`,`
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 z-index: 1;
 `),F(`content`,`
 font-size: var(--n-font-size);
 margin: var(--n-content-margin);
 position: relative;
 word-break: break-word;
 `,[L(`last`,`margin-bottom: 0;`)]),F(`action`,`
 display: flex;
 justify-content: flex-end;
 `,[C(`> *:not(:last-child)`,`
 margin-right: var(--n-action-space);
 `)]),F(`icon`,`
 font-size: var(--n-icon-size);
 transition: color .3s var(--n-bezier);
 `),F(`title`,`
 transition: color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 font-size: var(--n-title-font-size);
 font-weight: var(--n-title-font-weight);
 color: var(--n-title-text-color);
 `),b(`dialog-icon-container`,`
 display: flex;
 justify-content: center;
 `)]),m(b(`dialog`,`
 width: 446px;
 max-width: calc(100vw - 32px);
 `)),b(`dialog`,[k(`
 width: 446px;
 max-width: calc(100vw - 32px);
 `)])]),Ue={default:()=>(A(),p(je)),info:()=>(A(),p(je)),success:()=>(A(),p(Ee)),warning:()=>(A(),p(Me)),error:()=>(A(),p(De))},We=h({name:`Dialog`,alias:[`NimbusConfirmCard`,`Confirm`],props:{...D.props,...Z},slots:Object,setup(e){let{mergedComponentPropsRef:t,mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedRtlRef:a}=d(e),o=E(`Dialog`,a,n),s=i(()=>{let{iconPlacement:n}=e;return n||t?.value?.Dialog?.iconPlacement||`left`});function c(t){let{onPositiveClick:n}=e;n&&n(t)}function l(t){let{onNegativeClick:n}=e;n&&n(t)}function u(){let{onClose:t}=e;t&&t()}let f=D(`Dialog`,`-dialog`,He,Ae,e,n),p=i(()=>{let{type:t}=e,n=s.value,{common:{cubicBezierEaseInOut:r},self:{fontSize:i,lineHeight:a,border:o,titleTextColor:c,textColor:l,color:u,closeBorderRadius:d,closeColorHover:p,closeColorPressed:m,closeIconColor:h,closeIconColorHover:g,closeIconColorPressed:_,closeIconSize:v,borderRadius:y,titleFontWeight:b,titleFontSize:x,padding:S,iconSize:C,actionSpace:w,contentMargin:T,closeSize:E,[n===`top`?`iconMarginIconTop`:`iconMargin`]:D,[n===`top`?`closeMarginIconTop`:`closeMargin`]:O,[ne(`iconColor`,t)]:k}}=f.value,A=oe(D);return{"--n-font-size":i,"--n-icon-color":k,"--n-bezier":r,"--n-close-margin":O,"--n-icon-margin-top":A.top,"--n-icon-margin-right":A.right,"--n-icon-margin-bottom":A.bottom,"--n-icon-margin-left":A.left,"--n-icon-size":C,"--n-close-size":E,"--n-close-icon-size":v,"--n-close-border-radius":d,"--n-close-color-hover":p,"--n-close-color-pressed":m,"--n-close-icon-color":h,"--n-close-icon-color-hover":g,"--n-close-icon-color-pressed":_,"--n-color":u,"--n-text-color":l,"--n-border-radius":y,"--n-padding":S,"--n-line-height":a,"--n-border":o,"--n-content-margin":T,"--n-title-font-size":x,"--n-title-font-weight":b,"--n-title-text-color":c,"--n-action-space":w}}),m=r?x(`dialog`,i(()=>`${e.type[0]}${s.value[0]}`),p,e):void 0;return{mergedClsPrefix:n,rtlEnabled:o,mergedIconPlacement:s,mergedTheme:f,handlePositiveClick:c,handleNegativeClick:l,handleCloseClick:u,cssVars:r?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){let{bordered:n,mergedIconPlacement:r,cssVars:i,closable:o,showIcon:s,title:u,content:d,action:f,negativeText:m,positiveText:h,positiveButtonProps:g,negativeButtonProps:_,handlePositiveClick:v,handleNegativeClick:b,mergedTheme:x,loading:S,type:C,mergedClsPrefix:T}=this;this.onRender?.();let E=s?(A(),p(y,{key:1,clsPrefix:T,class:t(`${T}-dialog__icon`)},{default:()=>z(this.$slots.icon,e=>e||(this.icon?q(this.icon):Ue[this.type]()))},1032,[`clsPrefix`,`class`])):null,D=z(this.$slots.action,n=>n||h||m||f?(A(),e(`div`,{key:2,class:t([`${T}-dialog__action`,this.actionClass]),style:I(this.actionStyle)},[w(()=>n||(f?[q(f)]:[this.negativeText&&(A(),p(W,a({key:3,theme:x.peers.Button,themeOverrides:x.peerOverrides.Button,ghost:!0,size:`small`,onClick:b},_),{default:()=>q(this.negativeText)},1040,[`theme`,`themeOverrides`,`onClick`])),this.positiveText&&(A(),p(W,a({key:4,theme:x.peers.Button,themeOverrides:x.peerOverrides.Button,size:`small`,type:C==="default"?`primary`:C,disabled:S,loading:S,onClick:v},g),{default:()=>q(this.positiveText)},1040,[`theme`,`themeOverrides`,`type`,`disabled`,`loading`,`onClick`]))]))],6)):null);return A(),e(`div`,{class:t([`${T}-dialog`,this.themeClass,this.closable&&`${T}-dialog--closable`,`${T}-dialog--icon-${r}`,n&&`${T}-dialog--bordered`,this.rtlEnabled&&`${T}-dialog--rtl`]),style:I(i),role:`dialog`},[o?(A(),e(l,{key:0},[w(()=>z(this.$slots.close,n=>{let r=[`${T}-dialog__close`,this.rtlEnabled&&`${T}-dialog--rtl`];return n?(A(),e(`div`,{key:5,class:t(r)},[w(()=>n)],2)):(A(),p(Oe,{key:6,focusable:this.closeFocusable,clsPrefix:T,class:t(r),onClick:this.handleCloseClick},null,8,[`focusable`,`clsPrefix`,`class`,`onClick`]))}))],64)):w(()=>null),s&&r===`top`?(A(),e(`div`,{key:2,class:t(`${T}-dialog-icon-container`)},[w(()=>E)],2)):w(()=>null),c(`div`,{class:t([`${T}-dialog__title`,this.titleClass]),style:I(this.titleStyle)},[s&&r===`left`?(A(),e(l,{key:0},[w(()=>E)],64)):w(()=>null),w(()=>H(this.$slots.header,()=>[q(u)]))],6),c(`div`,{class:t([`${T}-dialog__content`,D?``:`${T}-dialog__content--last`,this.contentClass]),style:I(this.contentStyle)},[w(()=>H(this.$slots.default,()=>[q(d)]))],6),w(()=>D)],6)}}),Q=`n-draggable`;function Ge(e,t){let n,r=S(null),a=S(null),o=i(()=>e.value!==!1),s=i(()=>o.value?Q:``),c=i(()=>{let t=e.value;return t===!0||t===!1||!t||t.bounds!==`none`});function l(e){let i=e.querySelector(`.${Q}`);if(!i||!s.value)return;let o=0,l=0,u=0,d=0,f=0,p=0,m,h=null,_=null;function v(t){t.preventDefault(),m=t;let{x:n,y:i,right:s,bottom:c}=e.getBoundingClientRect();if(l=n,d=i,o=window.innerWidth-s,u=window.innerHeight-c,r.value!==null&&a.value!==null)p=r.value,f=a.value;else{let{left:t,top:n}=e.style;f=+n.slice(0,-2),p=+t.slice(0,-2)}}function y(){_&&=(r.value=_.x,a.value=_.y,null),h=null}function b(e){if(!m)return;let{clientX:t,clientY:n}=m,r=e.clientX-t,i=e.clientY-n;c.value&&(r>o?r=o:-r>l&&(r=-l),i>u?i=u:-i>d&&(i=-d)),_={x:r+p,y:i+f},h||=requestAnimationFrame(y)}function x(){m=void 0,h&&=(cancelAnimationFrame(h),null),_&&=(r.value=_.x,a.value=_.y,null),g(()=>{t.onEnd(e)})}R(`mousedown`,i,v),R(`mousemove`,window,b),R(`mouseup`,window,x),n=()=>{h&&cancelAnimationFrame(h),U(`mousedown`,i,v),U(`mousemove`,window,b),U(`mouseup`,window,x)}}function u(){n&&=(n(),void 0),r.value=null,a.value=null}return ee(u),{stopDrag:u,startDrag:l,draggableRef:o,draggableClassRef:s,dragX:r,dragY:a}}var $={...we,...Z},Ke=B($).filter(e=>e!==`onClose`&&e!==`onPositiveClick`&&e!==`onNegativeClick`),qe=h({name:`ModalBody`,inheritAttrs:!1,slots:Object,props:{show:{type:Boolean,required:!0},preset:String,displayDirective:{type:String,required:!0},trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},blockScroll:Boolean,draggable:{type:[Boolean,Object],default:!1},maskHidden:Boolean,...$,onClickoutside:{type:Function,required:!0},onBeforeLeave:{type:Function,required:!0},onAfterLeave:{type:Function,required:!0},onPositiveClick:{type:Function,required:!0},onNegativeClick:{type:Function,required:!0},onClose:{type:Function,required:!0},onAfterEnter:Function,onEsc:Function},setup(e){let t=S(null),n=S(null),r=S(e.show),a=S(null),s=S(null),c=o(G),l=null;M(P(e,`show`),e=>{e&&(l=c.getMousePosition())},{immediate:!0});let{stopDrag:u,startDrag:d,draggableRef:f,draggableClassRef:p,dragX:m,dragY:h}=Ge(P(e,`draggable`),{onEnd:e=>{x(e)}}),_=i(()=>v([e.titleClass,p.value])),y=i(()=>v([e.headerClass,p.value]));M(P(e,`show`),e=>{e&&(r.value=!0)}),xe(i(()=>e.blockScroll&&r.value));function b(){if(c.transformOriginRef.value===`center`)return``;let{value:e}=a,{value:t}=s;return e===null||t===null?``:n.value?`${e}px ${t+n.value.containerScrollTop}px`:``}function x(e){if(c.transformOriginRef.value===`center`||!l||!n.value)return;let t=n.value.containerScrollTop,{offsetLeft:r,offsetTop:i}=e,o=l.y,u=l.x;a.value=-(r-u),s.value=-(i-o-t),e.style.transformOrigin=b()}function C(e){g(()=>{x(e)})}function w(t){t.style.transformOrigin=b(),e.onBeforeLeave()}function T(t){let n=t;f.value&&d(n),e.onAfterEnter&&e.onAfterEnter(n)}function E(){r.value=!1,a.value=null,s.value=null,u(),e.onAfterLeave()}function D(){let{onClose:t}=e;t&&t()}function O(){e.onNegativeClick()}function k(){e.onPositiveClick()}let A=S(null);return M(A,e=>{e&&g(()=>{let n=e.el;n&&t.value!==n&&(t.value=n)})}),j(ue,t),j(be,null),j(de,null),{mergedTheme:c.mergedThemeRef,appear:c.appearRef,isMounted:c.isMountedRef,mergedClsPrefix:c.mergedClsPrefixRef,bodyRef:t,scrollbarRef:n,draggableClass:p,displayed:r,childNodeRef:A,cardHeaderClass:y,dialogTitleClass:_,handlePositiveClick:k,handleNegativeClick:O,handleCloseClick:D,handleAfterEnter:T,handleAfterLeave:E,handleBeforeLeave:w,handleEnter:C,dragX:m,dragY:h}},render(){let{$slots:i,$attrs:o,handleEnter:s,handleAfterEnter:c,handleAfterLeave:l,handleBeforeLeave:d,preset:f,mergedClsPrefix:m,dragX:h,dragY:g}=this,_={...o};h!==null&&g!==null&&(_.style=I([_.style,{left:`${h}px`,top:`${g}px`}]));let v=null;if(!f){if(v=ve(`default`,i.default,{draggableClass:this.draggableClass}),!v){re(`modal`,`default slot is empty`);return}v=r(v),v.props=a({class:`${m}-modal`},_,v.props||{})}return this.displayDirective===`show`||this.displayed||this.show?N((A(),e(`div`,{key:1,role:`none`,class:t([`${m}-modal-body-wrapper`,this.maskHidden&&`${m}-modal-body-wrapper--mask-hidden`])},[(A(),p(ce,{ref:`scrollbarRef`,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:`${m}-modal-scroll-content`},{default:()=>(A(),p(me,{disabled:!this.trapFocus||this.maskHidden,active:this.show,onEsc:this.onEsc,autoFocus:this.autoFocus},{default:()=>(A(),p(n,{name:`fade-in-scale-up-transition`,appear:this.appear??this.isMounted,onEnter:s,onAfterEnter:c,onAfterLeave:l,onBeforeLeave:d},{default:()=>{let e=[[T,this.show]];return e.push([le,this.onClickoutside,void 0,{capture:!0}]),N(this.preset===`confirm`||this.preset===`dialog`?(A(),p(We,a({key:2},_,{class:[`${m}-modal`,_.class],ref:`bodyRef`,theme:this.mergedTheme.peers.Dialog,themeOverrides:this.mergedTheme.peerOverrides.Dialog},K(this.$props,Ve),{titleClass:this.dialogTitleClass,"aria-modal":`true`}),u(i),1040,[`class`,`theme`,`themeOverrides`,`titleClass`])):this.preset===`card`?(A(),p(Te,a({key:3},_,{ref:`bodyRef`,class:[`${m}-modal`,_.class],theme:this.mergedTheme.peers.Card,themeOverrides:this.mergedTheme.peerOverrides.Card},K(this.$props,Ce),{headerClass:this.cardHeaderClass,"aria-modal":`true`,role:`dialog`}),u(i),1040,[`class`,`theme`,`themeOverrides`,`headerClass`])):this.childNodeRef=v,e)}},1032,[`appear`,`onEnter`,`onAfterEnter`,`onAfterLeave`,`onBeforeLeave`]))},1032,[`disabled`,`active`,`onEsc`,`autoFocus`]))},1032,[`theme`,`themeOverrides`,`contentClass`]))],2)),[[T,this.displayDirective===`if`||this.displayed||this.show]]):null}}),Je=C([b(`modal-container`,`
 position: fixed;
 left: 0;
 top: 0;
 height: 0;
 width: 0;
 display: flex;
 `),b(`modal-mask`,`
 position: fixed;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 background-color: rgba(0, 0, 0, .4);
 `,[se({enterDuration:`.25s`,leaveDuration:`.25s`,enterCubicBezier:`var(--n-bezier-ease-out)`,leaveCubicBezier:`var(--n-bezier-ease-out)`})]),b(`modal-body-wrapper`,`
 position: fixed;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 overflow: visible;
 `,[b(`modal-scroll-content`,`
 min-height: 100%;
 display: flex;
 position: relative;
 `),L(`mask-hidden`,`pointer-events: none;`,[b(`modal-scroll-content`,[C(`> *`,`
 pointer-events: all;
 `)])])]),b(`modal`,`
 position: relative;
 align-self: center;
 color: var(--n-text-color);
 margin: auto;
 box-shadow: var(--n-box-shadow);
 `,[_e({duration:`.25s`,enterScale:`.5`}),C(`.${Q}`,`
 cursor: move;
 user-select: none;
 `)])]),Ye={...D.props,show:Boolean,showMask:{type:Boolean,default:!0},maskClosable:{type:Boolean,default:!0},preset:String,to:[String,Object],displayDirective:{type:String,default:`if`},transformOrigin:{type:String,default:`mouse`},zIndex:Number,autoFocus:{type:Boolean,default:!0},trapFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},...$,draggable:[Boolean,Object],onEsc:Function,"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onBeforeLeave:Function,onAfterLeave:Function,onClose:Function,onPositiveClick:Function,onNegativeClick:Function,onMaskClick:Function,internalDialog:Boolean,internalModal:Boolean,internalAppear:{type:Boolean,default:void 0},overlayStyle:[String,Object],onBeforeHide:Function,onAfterHide:Function,onHide:Function,unstableShowMask:{type:Boolean,default:void 0}},Xe=h({name:`Modal`,inheritAttrs:!1,props:Ye,slots:Object,setup(e){let t=S(null),{mergedClsPrefixRef:n,namespaceRef:r,inlineThemeDisabled:a}=d(e),s=D(`Modal`,`-modal`,Je,ke,e,n),c=ze(64),l=Fe(),u=te(),f=e.internalDialog?o(Be,null):null,p=e.internalModal?o(ye,null):null,m=Se();function h(t){let{onUpdateShow:n,"onUpdate:show":r,onHide:i}=e;n&&V(n,t),r&&V(r,t),i&&!t&&i(t)}function g(){let{onClose:t}=e;t?Promise.resolve(t()).then(e=>{e!==!1&&h(!1)}):h(!1)}function _(){let{onPositiveClick:t}=e;t?Promise.resolve(t()).then(e=>{e!==!1&&h(!1)}):h(!1)}function v(){let{onNegativeClick:t}=e;t?Promise.resolve(t()).then(e=>{e!==!1&&h(!1)}):h(!1)}function y(){let{onBeforeLeave:t,onBeforeHide:n}=e;t&&V(t),n&&n()}function b(){let{onAfterLeave:t,onAfterHide:n}=e;t&&V(t),n&&n()}function C(n){let{onMaskClick:r}=e;r&&r(n),e.maskClosable&&t.value?.contains(ae(n))&&h(!1)}function w(t){e.onEsc?.(),e.show&&e.closeOnEsc&&ie(t)&&(m.value||h(!1))}j(G,{getMousePosition:()=>{let e=f||p;if(e){let{clickedRef:t,clickedPositionRef:n}=e;if(t.value&&n.value)return n.value}return c.value?l.value:null},mergedClsPrefixRef:n,mergedThemeRef:s,isMountedRef:u,appearRef:P(e,`internalAppear`),transformOriginRef:P(e,`transformOrigin`)});let T=i(()=>{let{common:{cubicBezierEaseOut:e},self:{boxShadow:t,color:n,textColor:r}}=s.value;return{"--n-bezier-ease-out":e,"--n-box-shadow":t,"--n-color":n,"--n-text-color":r}}),E=a?x(`theme-class`,void 0,T,e):void 0;return{mergedClsPrefix:n,namespace:r,isMounted:u,containerRef:t,presetProps:i(()=>K(e,Ke)),handleEsc:w,handleAfterLeave:b,handleClickoutside:C,handleBeforeLeave:y,doUpdateShow:h,handleNegativeClick:v,handlePositiveClick:_,handleCloseClick:g,cssVars:a?void 0:T,themeClass:E?.themeClass,onRender:E?.onRender}},render(){let{mergedClsPrefix:r}=this;return A(),p(ge,{to:this.to,show:this.show},{default:()=>{this.onRender?.();let{showMask:i}=this;return N((A(),e(`div`,{role:`none`,ref:`containerRef`,class:t([`${r}-modal-container`,this.themeClass,this.namespace]),style:I(this.cssVars)},[i?(A(),p(n,{name:`fade-in-transition`,key:`mask`,appear:this.internalAppear??this.isMounted},{default:()=>this.show?(A(),e(`div`,{key:1,"aria-hidden":!0,class:t(`${r}-modal-mask`)},null,2)):null},1032,[`appear`])):w(()=>null),(A(),p(qe,a({style:this.overlayStyle},this.$attrs,{ref:`bodyWrapper`,displayDirective:this.displayDirective,show:this.show,preset:this.preset,autoFocus:this.autoFocus,trapFocus:this.trapFocus,draggable:this.draggable,blockScroll:this.blockScroll,maskHidden:!i},this.presetProps,{onEsc:this.handleEsc,onClose:this.handleCloseClick,onNegativeClick:this.handleNegativeClick,onPositiveClick:this.handlePositiveClick,onBeforeLeave:this.handleBeforeLeave,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave,onClickoutside:this.handleClickoutside}),u(this.$slots),1040,[`style`,`displayDirective`,`show`,`preset`,`autoFocus`,`trapFocus`,`draggable`,`blockScroll`,`maskHidden`,`onEsc`,`onClose`,`onNegativeClick`,`onPositiveClick`,`onBeforeLeave`,`onAfterEnter`,`onAfterLeave`,`onClickoutside`]))],6)),[[fe,{zIndex:this.zIndex,enabled:this.show}]])}},1032,[`to`,`show`])}});export{Xe as t};