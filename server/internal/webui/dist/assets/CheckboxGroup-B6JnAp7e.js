import{C as e,F as t,P as n,X as r,d as i,et as a,f as o,h as s,o as c,p as l,v as u,vt as d,y as f}from"./vue-i18n-CU1juWHN.js";import{Ct as p,E as m,L as h,P as g,St as _,T as v,a as y,d as b,gt as x,h as S,ht as C,i as w,k as T,l as E,m as D,vt as O,xt as k,y as A,yt as j}from"./light-uQ0rL05w.js";import{v as M}from"./event-D2J5m536.js";import{F as N,L as P,P as F,x as I}from"./http-LAjieWQw.js";import{w as L}from"./index-D1aRHhil.js";var R=()=>(()=>{let e=v(`75be776d8875fa17`);return e[0]||=o(`svg`,{viewBox:`0 0 64 64`,class:`check-icon`},[o(`path`,{d:`M50.42,16.76L22.34,39.45l-8.1-11.46c-1.12-1.58-3.3-1.96-4.88-0.84c-1.58,1.12-1.95,3.3-0.84,4.88l10.26,14.51  c0.56,0.79,1.42,1.31,2.38,1.45c0.16,0.02,0.32,0.03,0.48,0.03c0.8,0,1.57-0.27,2.2-0.78l30.99-25.03c1.5-1.21,1.74-3.42,0.52-4.92  C54.13,15.78,51.93,15.55,50.42,16.76z`})],-1)})(),z=()=>(()=>{let e=v(`c6eed899356c8404`);return e[0]||=o(`svg`,{viewBox:`0 0 100 100`,class:`line-icon`},[o(`path`,{d:`M80.2,55.5H21.4c-2.8,0-5.1-2.5-5.1-5.5l0,0c0-3,2.3-5.5,5.1-5.5h58.7c2.8,0,5.1,2.5,5.1,5.5l0,0C85.2,53.1,82.9,55.5,80.2,55.5z`})],-1)})(),B=C([x(`checkbox`,`
 font-size: var(--n-font-size);
 outline: none;
 cursor: pointer;
 display: inline-flex;
 flex-wrap: nowrap;
 align-items: flex-start;
 word-break: break-word;
 line-height: var(--n-size);
 --n-merged-color-table: var(--n-color-table);
 `,[j(`show-label`,`line-height: var(--n-label-line-height);`),C(`&:hover`,[x(`checkbox-box`,[O(`border`,`border: var(--n-border-checked);`)])]),C(`&:focus:not(:active)`,[x(`checkbox-box`,[O(`border`,`
 border: var(--n-border-focus);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),j(`inside-table`,[x(`checkbox-box`,`
 background-color: var(--n-merged-color-table);
 `)]),j(`checked`,[x(`checkbox-box`,`
 background-color: var(--n-color-checked);
 `,[x(`checkbox-icon`,[C(`.check-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),j(`indeterminate`,[x(`checkbox-box`,[x(`checkbox-icon`,[C(`.check-icon`,`
 opacity: 0;
 transform: scale(.5);
 `),C(`.line-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),j(`checked, indeterminate`,[C(`&:focus:not(:active)`,[x(`checkbox-box`,[O(`border`,`
 border: var(--n-border-checked);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),x(`checkbox-box`,`
 background-color: var(--n-color-checked);
 border-left: 0;
 border-top: 0;
 `,[O(`border`,{border:`var(--n-border-checked)`})])]),j(`disabled`,{cursor:`not-allowed`},[j(`checked`,[x(`checkbox-box`,`
 background-color: var(--n-color-disabled-checked);
 `,[O(`border`,{border:`var(--n-border-disabled-checked)`}),x(`checkbox-icon`,[C(`.check-icon, .line-icon`,{fill:`var(--n-check-mark-color-disabled-checked)`})])])]),x(`checkbox-box`,`
 background-color: var(--n-color-disabled);
 `,[O(`border`,`
 border: var(--n-border-disabled);
 `),x(`checkbox-icon`,[C(`.check-icon, .line-icon`,`
 fill: var(--n-check-mark-color-disabled);
 `)])]),O(`label`,`
 color: var(--n-text-color-disabled);
 `)]),x(`checkbox-box-wrapper`,`
 position: relative;
 width: var(--n-size);
 flex-shrink: 0;
 flex-grow: 0;
 user-select: none;
 -webkit-user-select: none;
 `),x(`checkbox-box`,`
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
 `,[O(`border`,`
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
 `),x(`checkbox-icon`,`
 display: flex;
 align-items: center;
 justify-content: center;
 position: absolute;
 left: 1px;
 right: 1px;
 top: 1px;
 bottom: 1px;
 `,[C(`.check-icon, .line-icon`,`
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
 `),w({left:`1px`,top:`1px`})])]),O(`label`,`
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 `,[C(`&:empty`,{display:`none`})])]),_(x(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-modal);
 `)),p(x(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-popover);
 `))]),V=[`id`],H=[`tabindex`,`aria-checked`,`aria-labelledby`,`onKeyup`,`onKeydown`,`onClick`],U={...D.props,size:String,checked:{type:[Boolean,String,Number],default:void 0},defaultChecked:{type:[Boolean,String,Number],default:!1},value:[String,Number],disabled:{type:Boolean,default:void 0},indeterminate:Boolean,label:String,focusable:{type:Boolean,default:!0},checkedValue:{type:[Boolean,String,Number],default:!0},uncheckedValue:{type:[Boolean,String,Number],default:!1},"onUpdate:checked":[Function,Array],onUpdateChecked:[Function,Array],privateInsideTable:Boolean,onChange:[Function,Array]},W=f({name:`Checkbox`,props:U,setup(t){let n=e(G,null),o=r(null),{mergedClsPrefixRef:s,inlineThemeDisabled:c,mergedRtlRef:l,mergedComponentPropsRef:u}=g(t),d=r(t.defaultChecked),f=a(t,`checked`),p=M(f,d),m=b(()=>{if(n){let e=n.valueSetRef.value;return e&&t.value!==void 0?e.has(t.value):!1}return p.value===t.checkedValue}),h=I(t,{mergedSize(e){let{size:r}=t;if(r!==void 0)return r;if(n){let{value:e}=n.mergedSizeRef;if(e!==void 0)return e}if(e){let{mergedSize:t}=e;if(t!==void 0)return t.value}return u?.value?.Checkbox?.size||`medium`},mergedDisabled(e){let{disabled:r}=t;if(r!==void 0)return r;if(n){if(n.disabledRef.value)return!0;let{maxRef:{value:e},checkedCountRef:t}=n;if(e!==void 0&&t.value>=e&&!m.value)return!0;let{minRef:{value:r}}=n;if(r!==void 0&&t.value<=r&&m.value)return!0}return e?e.disabled.value:!1}}),{mergedDisabledRef:_,mergedSizeRef:v}=h,y=D(`Checkbox`,`-checkbox`,B,L,t,s);function x(e){if(n&&t.value!==void 0)n.toggleCheckbox(!m.value,t.value);else{let{onChange:n,"onUpdate:checked":r,onUpdateChecked:i}=t,{nTriggerFormInput:a,nTriggerFormChange:o}=h,s=m.value?t.uncheckedValue:t.checkedValue;r&&N(r,s,e),i&&N(i,s,e),n&&N(n,s,e),a(),o(),d.value=s}}function C(e){_.value||x(e)}function w(e){if(!_.value)switch(e.key){case` `:case`Enter`:x(e)}}function T(e){e.key===` `&&e.preventDefault()}let O={focus:()=>{o.value?.focus()},blur:()=>{o.value?.blur()}},j=E(`Checkbox`,l,s),P=i(()=>{let{value:e}=v,{common:{cubicBezierEaseInOut:t},self:{borderRadius:n,color:r,colorChecked:i,colorDisabled:a,colorTableHeader:o,colorTableHeaderModal:s,colorTableHeaderPopover:c,checkMarkColor:l,checkMarkColorDisabled:u,border:d,borderFocus:f,borderDisabled:p,borderChecked:m,boxShadowFocus:h,textColor:g,textColorDisabled:_,checkMarkColorDisabledChecked:b,colorDisabledChecked:x,borderDisabledChecked:S,labelPadding:C,labelLineHeight:w,labelFontWeight:T,[k(`fontSize`,e)]:E,[k(`size`,e)]:D}}=y.value;return{"--n-label-line-height":w,"--n-label-font-weight":T,"--n-size":D,"--n-bezier":t,"--n-border-radius":n,"--n-border":d,"--n-border-checked":m,"--n-border-focus":f,"--n-border-disabled":p,"--n-border-disabled-checked":S,"--n-box-shadow-focus":h,"--n-color":r,"--n-color-checked":i,"--n-color-table":o,"--n-color-table-modal":s,"--n-color-table-popover":c,"--n-color-disabled":a,"--n-color-disabled-checked":x,"--n-text-color":g,"--n-text-color-disabled":_,"--n-check-mark-color":l,"--n-check-mark-color-disabled":u,"--n-check-mark-color-disabled-checked":b,"--n-font-size":E,"--n-label-padding":C}}),F=c?S(`checkbox`,i(()=>v.value[0]),P,t):void 0;return Object.assign(h,O,{rtlEnabled:j,selfRef:o,mergedClsPrefix:s,mergedDisabled:_,renderedChecked:m,mergedTheme:y,labelId:A(),handleClick:C,handleKeyUp:w,handleKeyDown:T,cssVars:c?void 0:P,themeClass:F?.themeClass,onRender:F?.onRender})},render(){let{$slots:e,renderedChecked:t,mergedDisabled:r,indeterminate:i,privateInsideTable:a,cssVars:c,labelId:l,label:f,mergedClsPrefix:p,focusable:h,handleKeyUp:g,handleKeyDown:_,handleClick:b}=this;this.onRender?.();let x=F(e.default,e=>f||e?(n(),s(`span`,{key:1,class:m(`${p}-checkbox__label`),id:l},[T(()=>f||e)],10,V)):null);return(()=>{let e=v(`70be6e74cd27cb50`);return n(),s(`div`,{ref:`selfRef`,class:m([`${p}-checkbox`,this.themeClass,this.rtlEnabled&&`${p}-checkbox--rtl`,t&&`${p}-checkbox--checked`,r&&`${p}-checkbox--disabled`,i&&`${p}-checkbox--indeterminate`,a&&`${p}-checkbox--inside-table`,x&&`${p}-checkbox--show-label`]),tabindex:r||!h?void 0:0,role:`checkbox`,"aria-checked":i?`mixed`:t,"aria-labelledby":l,style:d(c),onKeyup:g,onKeydown:_,onClick:b,onMousedown:e[0]||=()=>{P(`selectstart`,window,e=>{e.preventDefault()},{once:!0})}},[o(`div`,{class:m(`${p}-checkbox-box-wrapper`)},[e[1]||=T(`\xA0`,-1),o(`div`,{class:m(`${p}-checkbox-box`)},[u(y,null,{default:()=>this.indeterminate?(n(),s(`div`,{key:`indeterminate`,class:m(`${p}-checkbox-icon`)},[T(()=>z())],2)):(n(),s(`div`,{key:`check`,class:m(`${p}-checkbox-icon`)},[T(()=>R())],2))},1024),o(`div`,{class:m(`${p}-checkbox-box__border`)},null,2)],2)],2),T(()=>x)],46,H)})()}}),G=h(`n-checkbox-group`),K=f({name:`CheckboxGroup`,props:{min:Number,max:Number,size:String,options:Array,labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},value:Array,defaultValue:{type:Array,default:null},disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]},setup(e){let{mergedClsPrefixRef:n}=g(e),o=I(e),{mergedSizeRef:s,mergedDisabledRef:c}=o,l=r(e.defaultValue),u=i(()=>e.value),d=M(u,l),f=i(()=>d.value?.length||0),p=i(()=>Array.isArray(d.value)?new Set(d.value):new Set);function m(t,n){let{nTriggerFormInput:r,nTriggerFormChange:i}=o,{onChange:a,"onUpdate:value":s,onUpdateValue:c}=e;if(Array.isArray(d.value)){let e=Array.from(d.value),o=e.findIndex(e=>e===n);t?~o||(e.push(n),c&&N(c,e,{actionType:`check`,value:n}),s&&N(s,e,{actionType:`check`,value:n}),r(),i(),l.value=e,a&&N(a,e)):~o&&(e.splice(o,1),c&&N(c,e,{actionType:`uncheck`,value:n}),s&&N(s,e,{actionType:`uncheck`,value:n}),a&&N(a,e),l.value=e,r(),i())}else t?(c&&N(c,[n],{actionType:`check`,value:n}),s&&N(s,[n],{actionType:`check`,value:n}),a&&N(a,[n]),l.value=[n],r(),i()):(c&&N(c,[],{actionType:`uncheck`,value:n}),s&&N(s,[],{actionType:`uncheck`,value:n}),a&&N(a,[]),l.value=[],r(),i())}return t(G,{checkedCountRef:f,maxRef:a(e,`max`),minRef:a(e,`min`),valueSetRef:p,disabledRef:c,mergedSizeRef:s,toggleCheckbox:m}),{mergedClsPrefix:n}},render(){let{options:e,labelField:t,valueField:r}=this.$props;return n(),s(`div`,{class:m(`${this.mergedClsPrefix}-checkbox-group`),role:`group`},[e?(n(),s(c,{key:0},[T(()=>e.map(e=>{let i=e[r];return n(),l(W,{key:i,value:i,disabled:e.disabled,label:e[t]},null,8,[`value`,`disabled`,`label`])}))],64)):(n(),s(c,{key:1},[T(()=>this.$slots.default?.())],64))],2)}});export{W as n,K as t};