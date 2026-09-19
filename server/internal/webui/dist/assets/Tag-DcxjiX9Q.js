import{At as e,D as t,Ft as n,Jt as r,Mt as i,Ot as a,Yt as o,at as s,ct as c,d as l,dt as u,g as d,kt as f,lt as p,o as m,ot as h,pn as g,r as _,s as v,sn as y,u as b,un as x,ut as S,w as C,y as w}from"./vue-i18n-CgKceEWx.js";import{m as T,p as E,y as D}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{t as O}from"./color-to-class-G4jmrDIg.js";import{L as k,z as A}from"./index-kjlF3dpY.js";function j(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:c,baseColor:u,borderColor:d,opacityDisabled:f,tagColor:p,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:j}=e;return{...A,closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:f,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:j,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:u,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${d}`,textColor:t,color:p,colorBordered:`rgb(250, 250, 252)`,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${l(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:l(i,{alpha:.12}),colorBorderedPrimary:l(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:l(i,{alpha:.12}),closeColorPressedPrimary:l(i,{alpha:.18}),borderInfo:`1px solid ${l(a,{alpha:.3})}`,textColorInfo:a,colorInfo:l(a,{alpha:.12}),colorBorderedInfo:l(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:l(a,{alpha:.12}),closeColorPressedInfo:l(a,{alpha:.18}),borderSuccess:`1px solid ${l(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:l(o,{alpha:.12}),colorBorderedSuccess:l(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:l(o,{alpha:.12}),closeColorPressedSuccess:l(o,{alpha:.18}),borderWarning:`1px solid ${l(s,{alpha:.35})}`,textColorWarning:s,colorWarning:l(s,{alpha:.15}),colorBorderedWarning:l(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:l(s,{alpha:.12}),closeColorPressedWarning:l(s,{alpha:.18}),borderError:`1px solid ${l(c,{alpha:.23})}`,textColorError:c,colorError:l(c,{alpha:.1}),colorBorderedError:l(c,{alpha:.08}),closeIconColorError:c,closeIconColorHoverError:c,closeIconColorPressedError:c,closeColorHoverError:l(c,{alpha:.12}),closeColorPressedError:l(c,{alpha:.18})}}var M={name:`Tag`,common:b,self:j},N={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},P=h(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[p(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),c(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),c(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),c(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),c(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),p(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[c(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),c(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),p(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),p(`icon, avatar`,[p(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),p(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),p(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[S(`disabled`,[s(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[S(`checked`,`color: var(--n-text-color-hover-checkable);`)]),s(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[S(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),p(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[S(`disabled`,[s(`&:hover`,`background-color: var(--n-color-checked-hover);`),s(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),F=[`onClick`,`onMouseenter`,`onMouseleave`],I={...m.props,...N,bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function},L=t(`n-tag`),R=n({name:`Tag`,props:I,slots:Object,setup(e){let t=y(null),{mergedBorderedRef:n,mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedRtlRef:s,mergedComponentPropsRef:c}=C(e),l=a(()=>e.size||c?.value?.Tag?.size||`medium`),d=m(`Tag`,`-tag`,P,M,e,r);o(L,{roundRef:x(e,`round`)});function f(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function p(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&T(n,t)}}let h={setTextContent(e){let{value:n}=t;n&&(n.textContent=e)}},g=_(`Tag`,s,r),b=a(()=>{let{type:t,color:{color:r,textColor:i}={}}=e,a=l.value,{common:{cubicBezierEaseInOut:o},self:{padding:s,closeMargin:c,borderRadius:f,opacityDisabled:p,textColorCheckable:m,textColorHoverCheckable:h,textColorPressedCheckable:g,textColorChecked:_,colorCheckable:v,colorHoverCheckable:y,colorPressedCheckable:b,colorChecked:x,colorCheckedHover:S,colorCheckedPressed:C,closeBorderRadius:w,fontWeightStrong:T,[u(`colorBordered`,t)]:E,[u(`closeSize`,a)]:O,[u(`closeIconSize`,a)]:k,[u(`fontSize`,a)]:A,[u(`height`,a)]:j,[u(`color`,t)]:M,[u(`textColor`,t)]:N,[u(`border`,t)]:P,[u(`closeIconColor`,t)]:F,[u(`closeIconColorHover`,t)]:I,[u(`closeIconColorPressed`,t)]:L,[u(`closeColorHover`,t)]:R,[u(`closeColorPressed`,t)]:z}}=d.value,B=D(c);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${j} - 8px)`,"--n-bezier":o,"--n-border-radius":f,"--n-border":P,"--n-close-icon-size":k,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":w,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":O,"--n-color":r||(n.value?E:M),"--n-color-checkable":v,"--n-color-checked":x,"--n-color-checked-hover":S,"--n-color-checked-pressed":C,"--n-color-hover-checkable":y,"--n-color-pressed-checkable":b,"--n-font-size":A,"--n-height":j,"--n-opacity-disabled":p,"--n-padding":s,"--n-text-color":i||N,"--n-text-color-checkable":m,"--n-text-color-checked":_,"--n-text-color-hover-checkable":h,"--n-text-color-pressed-checkable":g}}),S=i?v(`tag`,a(()=>{let t=``,{type:r,color:{color:i,textColor:a}={}}=e;return t+=r[0],t+=l.value[0],i&&(t+=`a${O(i)}`),a&&(t+=`b${O(a)}`),n.value&&(t+=`c`),t}),b,e):void 0;return{...h,rtlEnabled:g,mergedClsPrefix:r,contentRef:t,mergedBordered:n,handleClick:f,handleCloseClick:p,cssVars:i?void 0:b,themeClass:S?.themeClass,onRender:S?.onRender}},render(){let{mergedClsPrefix:t,rtlEnabled:n,closable:a,color:{borderColor:o}={},round:s,onRender:c,$slots:l}=this;c?.();let u=E(l.avatar,e=>e&&(r(),i(`div`,{class:d(`${t}-tag__avatar`)},[w(()=>e)],2))),p=E(l.icon,e=>e&&(r(),i(`div`,{class:d(`${t}-tag__icon`)},[w(()=>e)],2)));return r(),i(`div`,{class:d([`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:s,[`${t}-tag--avatar`]:u,[`${t}-tag--icon`]:p,[`${t}-tag--closable`]:a}]),style:g(this.cssVars),onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},[w(()=>p||u),f(`span`,{class:d(`${t}-tag__content`),ref:`contentRef`},[w(()=>this.$slots.default?.())],2),!this.checkable&&a?(r(),e(k,{key:0,clsPrefix:t,class:d(`${t}-tag__close`),disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:s,isButtonTag:this.internalCloseIsButtonTag,absolute:!0},null,8,[`clsPrefix`,`class`,`disabled`,`onClick`,`focusable`,`round`,`isButtonTag`])):w(()=>null),!this.checkable&&this.mergedBordered?(r(),i(`div`,{key:2,class:d(`${t}-tag__border`),style:g({borderColor:o})},null,6)):w(()=>null)],46,F)}});export{R as t};