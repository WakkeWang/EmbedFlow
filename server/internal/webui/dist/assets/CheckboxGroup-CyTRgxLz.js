import{Ct as e,E as t,L as n,Lt as r,Nt as i,P as a,Rt as o,Sn as s,St as c,T as l,Ut as u,Vt as d,Wt as f,d as p,gn as m,gt as h,h as g,ht as _,in as v,k as y,l as b,m as x,o as S,qt as C,rn as w,s as T,vt as E,xt as D,y as O,yn as k,yt as A,zt as j}from"./vue-i18n-pPSPCP6m.js";import{A as M,D as N,O as P,m as F}from"./_plugin-vue_export-helper-cUzcR9Kp.js";import{f as I}from"./fade-in-scale-up.cssr-BU07VhSR.js";import{w as L}from"./index-BNdp-2__.js";var R=()=>(()=>{let e=l(`75be776d8875fa17`);return e[0]||=o(`svg`,{viewBox:`0 0 64 64`,class:`check-icon`},[o(`path`,{d:`M50.42,16.76L22.34,39.45l-8.1-11.46c-1.12-1.58-3.3-1.96-4.88-0.84c-1.58,1.12-1.95,3.3-0.84,4.88l10.26,14.51  c0.56,0.79,1.42,1.31,2.38,1.45c0.16,0.02,0.32,0.03,0.48,0.03c0.8,0,1.57-0.27,2.2-0.78l30.99-25.03c1.5-1.21,1.74-3.42,0.52-4.92  C54.13,15.78,51.93,15.55,50.42,16.76z`})],-1)})(),z=()=>(()=>{let e=l(`c6eed899356c8404`);return e[0]||=o(`svg`,{viewBox:`0 0 100 100`,class:`line-icon`},[o(`path`,{d:`M80.2,55.5H21.4c-2.8,0-5.1-2.5-5.1-5.5l0,0c0-3,2.3-5.5,5.1-5.5h58.7c2.8,0,5.1,2.5,5.1,5.5l0,0C85.2,53.1,82.9,55.5,80.2,55.5z`})],-1)})(),B=_([h(`checkbox`,`
 font-size: var(--n-font-size);
 outline: none;
 cursor: pointer;
 display: inline-flex;
 flex-wrap: nowrap;
 align-items: flex-start;
 word-break: break-word;
 line-height: var(--n-size);
 --n-merged-color-table: var(--n-color-table);
 `,[A(`show-label`,`line-height: var(--n-label-line-height);`),_(`&:hover`,[h(`checkbox-box`,[E(`border`,`border: var(--n-border-checked);`)])]),_(`&:focus:not(:active)`,[h(`checkbox-box`,[E(`border`,`
 border: var(--n-border-focus);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),A(`inside-table`,[h(`checkbox-box`,`
 background-color: var(--n-merged-color-table);
 `)]),A(`checked`,[h(`checkbox-box`,`
 background-color: var(--n-color-checked);
 `,[h(`checkbox-icon`,[_(`.check-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),A(`indeterminate`,[h(`checkbox-box`,[h(`checkbox-icon`,[_(`.check-icon`,`
 opacity: 0;
 transform: scale(.5);
 `),_(`.line-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),A(`checked, indeterminate`,[_(`&:focus:not(:active)`,[h(`checkbox-box`,[E(`border`,`
 border: var(--n-border-checked);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),h(`checkbox-box`,`
 background-color: var(--n-color-checked);
 border-left: 0;
 border-top: 0;
 `,[E(`border`,{border:`var(--n-border-checked)`})])]),A(`disabled`,{cursor:`not-allowed`},[A(`checked`,[h(`checkbox-box`,`
 background-color: var(--n-color-disabled-checked);
 `,[E(`border`,{border:`var(--n-border-disabled-checked)`}),h(`checkbox-icon`,[_(`.check-icon, .line-icon`,{fill:`var(--n-check-mark-color-disabled-checked)`})])])]),h(`checkbox-box`,`
 background-color: var(--n-color-disabled);
 `,[E(`border`,`
 border: var(--n-border-disabled);
 `),h(`checkbox-icon`,[_(`.check-icon, .line-icon`,`
 fill: var(--n-check-mark-color-disabled);
 `)])]),E(`label`,`
 color: var(--n-text-color-disabled);
 `)]),h(`checkbox-box-wrapper`,`
 position: relative;
 width: var(--n-size);
 flex-shrink: 0;
 flex-grow: 0;
 user-select: none;
 -webkit-user-select: none;
 `),h(`checkbox-box`,`
 position: absolute;
 left: 0;
 top: 50%;
 transform: translateY(-50%);
 height: var(--n-size);
 width: var(--n-size);
 display: inline-block;
 box-sizing: border-box;
 border-radius: var(--n-border-radius);
 background-color: var(--n-color);
 transition: background-color 0.3s var(--n-bezier);
 `,[E(`border`,`
 transition:
 border-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border: var(--n-border);
 `),h(`checkbox-icon`,`
 display: flex;
 align-items: center;
 justify-content: center;
 position: absolute;
 left: 1px;
 right: 1px;
 top: 1px;
 bottom: 1px;
 `,[_(`.check-icon, .line-icon`,`
 width: 100%;
 fill: var(--n-check-mark-color);
 opacity: 0;
 transform: scale(0.5);
 transform-origin: center;
 transition:
 fill 0.3s var(--n-bezier),
 transform 0.3s var(--n-bezier),
 opacity 0.3s var(--n-bezier),
 border-color 0.3s var(--n-bezier);
 `),S({left:`1px`,top:`1px`})])]),E(`label`,`
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 `,[_(`&:empty`,{display:`none`})])]),c(h(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-modal);
 `)),e(h(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-popover);
 `))]),V=[`id`],H=[`tabindex`,`aria-checked`,`aria-labelledby`,`onKeyup`,`onKeydown`,`onClick`],U={...x.props,size:String,checked:{type:[Boolean,String,Number],default:void 0},defaultChecked:{type:[Boolean,String,Number],default:!1},value:[String,Number],disabled:{type:Boolean,default:void 0},indeterminate:Boolean,label:String,focusable:{type:Boolean,default:!0},checkedValue:{type:[Boolean,String,Number],default:!0},uncheckedValue:{type:[Boolean,String,Number],default:!1},"onUpdate:checked":[Function,Array],onUpdateChecked:[Function,Array],privateInsideTable:Boolean,onChange:[Function,Array]},W=f({name:`Checkbox`,props:U,setup(e){let t=C(G,null),n=m(null),{mergedClsPrefixRef:i,inlineThemeDisabled:o,mergedRtlRef:s,mergedComponentPropsRef:c}=a(e),l=m(e.defaultChecked),u=k(e,`checked`),d=I(u,l),f=p(()=>{if(t){let n=t.valueSetRef.value;return n&&e.value!==void 0?n.has(e.value):!1}return d.value===e.checkedValue}),h=F(e,{mergedSize(n){let{size:r}=e;if(r!==void 0)return r;if(t){let{value:e}=t.mergedSizeRef;if(e!==void 0)return e}if(n){let{mergedSize:e}=n;if(e!==void 0)return e.value}return c?.value?.Checkbox?.size||`medium`},mergedDisabled(n){let{disabled:r}=e;if(r!==void 0)return r;if(t){if(t.disabledRef.value)return!0;let{maxRef:{value:e},checkedCountRef:n}=t;if(e!==void 0&&n.value>=e&&!f.value)return!0;let{minRef:{value:r}}=t;if(r!==void 0&&n.value<=r&&f.value)return!0}return n?n.disabled.value:!1}}),{mergedDisabledRef:_,mergedSizeRef:v}=h,y=x(`Checkbox`,`-checkbox`,B,L,e,i);function S(n){if(t&&e.value!==void 0)t.toggleCheckbox(!f.value,e.value);else{let{onChange:t,"onUpdate:checked":r,onUpdateChecked:i}=e,{nTriggerFormInput:a,nTriggerFormChange:o}=h,s=f.value?e.uncheckedValue:e.checkedValue;r&&P(r,s,n),i&&P(i,s,n),t&&P(t,s,n),a(),o(),l.value=s}}function w(e){_.value||S(e)}function T(e){if(!_.value)switch(e.key){case` `:case`Enter`:S(e)}}function E(e){e.key===` `&&e.preventDefault()}let A={focus:()=>{n.value?.focus()},blur:()=>{n.value?.blur()}},j=b(`Checkbox`,s,i),M=r(()=>{let{value:e}=v,{common:{cubicBezierEaseInOut:t},self:{borderRadius:n,color:r,colorChecked:i,colorDisabled:a,colorTableHeader:o,colorTableHeaderModal:s,colorTableHeaderPopover:c,checkMarkColor:l,checkMarkColorDisabled:u,border:d,borderFocus:f,borderDisabled:p,borderChecked:m,boxShadowFocus:h,textColor:g,textColorDisabled:_,checkMarkColorDisabledChecked:b,colorDisabledChecked:x,borderDisabledChecked:S,labelPadding:C,labelLineHeight:w,labelFontWeight:T,[D(`fontSize`,e)]:E,[D(`size`,e)]:O}}=y.value;return{"--n-label-line-height":w,"--n-label-font-weight":T,"--n-size":O,"--n-bezier":t,"--n-border-radius":n,"--n-border":d,"--n-border-checked":m,"--n-border-focus":f,"--n-border-disabled":p,"--n-border-disabled-checked":S,"--n-box-shadow-focus":h,"--n-color":r,"--n-color-checked":i,"--n-color-table":o,"--n-color-table-modal":s,"--n-color-table-popover":c,"--n-color-disabled":a,"--n-color-disabled-checked":x,"--n-text-color":g,"--n-text-color-disabled":_,"--n-check-mark-color":l,"--n-check-mark-color-disabled":u,"--n-check-mark-color-disabled-checked":b,"--n-font-size":E,"--n-label-padding":C}}),N=o?g(`checkbox`,r(()=>v.value[0]),M,e):void 0;return Object.assign(h,A,{rtlEnabled:j,selfRef:n,mergedClsPrefix:i,mergedDisabled:_,renderedChecked:f,mergedTheme:y,labelId:O(),handleClick:w,handleKeyUp:T,handleKeyDown:E,cssVars:o?void 0:M,themeClass:N?.themeClass,onRender:N?.onRender})},render(){let{$slots:e,renderedChecked:n,mergedDisabled:r,indeterminate:i,privateInsideTable:a,cssVars:c,labelId:f,label:p,mergedClsPrefix:m,focusable:h,handleKeyUp:g,handleKeyDown:_,handleClick:v}=this;this.onRender?.();let b=N(e.default,e=>p||e?(w(),d(`span`,{key:1,class:t(`${m}-checkbox__label`),id:f},[y(()=>p||e)],10,V)):null);return(()=>{let e=l(`70be6e74cd27cb50`);return w(),d(`div`,{ref:`selfRef`,class:t([`${m}-checkbox`,this.themeClass,this.rtlEnabled&&`${m}-checkbox--rtl`,n&&`${m}-checkbox--checked`,r&&`${m}-checkbox--disabled`,i&&`${m}-checkbox--indeterminate`,a&&`${m}-checkbox--inside-table`,b&&`${m}-checkbox--show-label`]),tabindex:r||!h?void 0:0,role:`checkbox`,"aria-checked":i?`mixed`:n,"aria-labelledby":f,style:s(c),onKeyup:g,onKeydown:_,onClick:v,onMousedown:e[0]||=()=>{M(`selectstart`,window,e=>{e.preventDefault()},{once:!0})}},[o(`div`,{class:t(`${m}-checkbox-box-wrapper`)},[e[1]||=y(`\xA0`,-1),o(`div`,{class:t(`${m}-checkbox-box`)},[u(T,null,{default:()=>this.indeterminate?(w(),d(`div`,{key:`indeterminate`,class:t(`${m}-checkbox-icon`)},[y(()=>z())],2)):(w(),d(`div`,{key:`check`,class:t(`${m}-checkbox-icon`)},[y(()=>R())],2))},1024),o(`div`,{class:t(`${m}-checkbox-box__border`)},null,2)],2)],2),y(()=>b)],46,H)})()}}),G=n(`n-checkbox-group`),K=f({name:`CheckboxGroup`,props:{min:Number,max:Number,size:String,options:Array,labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},value:Array,defaultValue:{type:Array,default:null},disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]},setup(e){let{mergedClsPrefixRef:t}=a(e),n=F(e),{mergedSizeRef:i,mergedDisabledRef:o}=n,s=m(e.defaultValue),c=r(()=>e.value),l=I(c,s),u=r(()=>l.value?.length||0),d=r(()=>Array.isArray(l.value)?new Set(l.value):new Set);function f(t,r){let{nTriggerFormInput:i,nTriggerFormChange:a}=n,{onChange:o,"onUpdate:value":c,onUpdateValue:u}=e;if(Array.isArray(l.value)){let e=Array.from(l.value),n=e.findIndex(e=>e===r);t?~n||(e.push(r),u&&P(u,e,{actionType:`check`,value:r}),c&&P(c,e,{actionType:`check`,value:r}),i(),a(),s.value=e,o&&P(o,e)):~n&&(e.splice(n,1),u&&P(u,e,{actionType:`uncheck`,value:r}),c&&P(c,e,{actionType:`uncheck`,value:r}),o&&P(o,e),s.value=e,i(),a())}else t?(u&&P(u,[r],{actionType:`check`,value:r}),c&&P(c,[r],{actionType:`check`,value:r}),o&&P(o,[r]),s.value=[r],i(),a()):(u&&P(u,[],{actionType:`uncheck`,value:r}),c&&P(c,[],{actionType:`uncheck`,value:r}),o&&P(o,[]),s.value=[],i(),a())}return v(G,{checkedCountRef:u,maxRef:k(e,`max`),minRef:k(e,`min`),valueSetRef:d,disabledRef:o,mergedSizeRef:i,toggleCheckbox:f}),{mergedClsPrefix:t}},render(){let{options:e,labelField:n,valueField:r}=this.$props;return w(),d(`div`,{class:t(`${this.mergedClsPrefix}-checkbox-group`),role:`group`},[e?(w(),d(i,{key:0},[y(()=>e.map(e=>{let t=e[r];return w(),j(W,{key:t,value:t,disabled:e.disabled,label:e[n]},null,8,[`value`,`disabled`,`label`])}))],64)):(w(),d(i,{key:1},[y(()=>this.$slots.default?.())],64))],2)}});export{W as n,K as t};