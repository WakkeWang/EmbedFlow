import{At as e,D as t,Ft as n,Jt as r,Mt as i,Ot as a,Pt as o,Rt as s,Yt as c,at as l,ct as u,dt as d,ft as f,g as p,h as m,kt as h,lt as g,o as _,ot as v,pn as y,pt as b,r as x,s as S,sn as C,un as w,w as T,wt as E,y as D}from"./vue-i18n-CgKceEWx.js";import{g as O,m as k,p as A}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{f as j}from"./fade-in-scale-up.cssr-7akLf3Ad.js";import{i as M,r as N}from"./FadeInExpandTransition-D8IUT6y9.js";import{m as P}from"./http-BdILZoCd.js";import{A as F,j as I,w as L}from"./index-kjlF3dpY.js";var R=()=>(()=>{let e=m(`75be776d8875fa17`);return e[0]||=h(`svg`,{viewBox:`0 0 64 64`,class:`check-icon`},[h(`path`,{d:`M50.42,16.76L22.34,39.45l-8.1-11.46c-1.12-1.58-3.3-1.96-4.88-0.84c-1.58,1.12-1.95,3.3-0.84,4.88l10.26,14.51  c0.56,0.79,1.42,1.31,2.38,1.45c0.16,0.02,0.32,0.03,0.48,0.03c0.8,0,1.57-0.27,2.2-0.78l30.99-25.03c1.5-1.21,1.74-3.42,0.52-4.92  C54.13,15.78,51.93,15.55,50.42,16.76z`})],-1)})(),z=()=>(()=>{let e=m(`c6eed899356c8404`);return e[0]||=h(`svg`,{viewBox:`0 0 100 100`,class:`line-icon`},[h(`path`,{d:`M80.2,55.5H21.4c-2.8,0-5.1-2.5-5.1-5.5l0,0c0-3,2.3-5.5,5.1-5.5h58.7c2.8,0,5.1,2.5,5.1,5.5l0,0C85.2,53.1,82.9,55.5,80.2,55.5z`})],-1)})(),B=l([v(`checkbox`,`
 font-size: var(--n-font-size);
 outline: none;
 cursor: pointer;
 display: inline-flex;
 flex-wrap: nowrap;
 align-items: flex-start;
 word-break: break-word;
 line-height: var(--n-size);
 --n-merged-color-table: var(--n-color-table);
 `,[g(`show-label`,`line-height: var(--n-label-line-height);`),l(`&:hover`,[v(`checkbox-box`,[u(`border`,`border: var(--n-border-checked);`)])]),l(`&:focus:not(:active)`,[v(`checkbox-box`,[u(`border`,`
 border: var(--n-border-focus);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),g(`inside-table`,[v(`checkbox-box`,`
 background-color: var(--n-merged-color-table);
 `)]),g(`checked`,[v(`checkbox-box`,`
 background-color: var(--n-color-checked);
 `,[v(`checkbox-icon`,[l(`.check-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),g(`indeterminate`,[v(`checkbox-box`,[v(`checkbox-icon`,[l(`.check-icon`,`
 opacity: 0;
 transform: scale(.5);
 `),l(`.line-icon`,`
 opacity: 1;
 transform: scale(1);
 `)])])]),g(`checked, indeterminate`,[l(`&:focus:not(:active)`,[v(`checkbox-box`,[u(`border`,`
 border: var(--n-border-checked);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),v(`checkbox-box`,`
 background-color: var(--n-color-checked);
 border-left: 0;
 border-top: 0;
 `,[u(`border`,{border:`var(--n-border-checked)`})])]),g(`disabled`,{cursor:`not-allowed`},[g(`checked`,[v(`checkbox-box`,`
 background-color: var(--n-color-disabled-checked);
 `,[u(`border`,{border:`var(--n-border-disabled-checked)`}),v(`checkbox-icon`,[l(`.check-icon, .line-icon`,{fill:`var(--n-check-mark-color-disabled-checked)`})])])]),v(`checkbox-box`,`
 background-color: var(--n-color-disabled);
 `,[u(`border`,`
 border: var(--n-border-disabled);
 `),v(`checkbox-icon`,[l(`.check-icon, .line-icon`,`
 fill: var(--n-check-mark-color-disabled);
 `)])]),u(`label`,`
 color: var(--n-text-color-disabled);
 `)]),v(`checkbox-box-wrapper`,`
 position: relative;
 width: var(--n-size);
 flex-shrink: 0;
 flex-grow: 0;
 user-select: none;
 -webkit-user-select: none;
 `),v(`checkbox-box`,`
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
 `,[u(`border`,`
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
 `),v(`checkbox-icon`,`
 display: flex;
 align-items: center;
 justify-content: center;
 position: absolute;
 left: 1px;
 right: 1px;
 top: 1px;
 bottom: 1px;
 `,[l(`.check-icon, .line-icon`,`
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
 `),F({left:`1px`,top:`1px`})])]),u(`label`,`
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 `,[l(`&:empty`,{display:`none`})])]),f(v(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-modal);
 `)),b(v(`checkbox`,`
 --n-merged-color-table: var(--n-color-table-popover);
 `))]),V=[`id`],H=[`tabindex`,`aria-checked`,`aria-labelledby`,`onKeyup`,`onKeydown`,`onClick`],U={..._.props,size:String,checked:{type:[Boolean,String,Number],default:void 0},defaultChecked:{type:[Boolean,String,Number],default:!1},value:[String,Number],disabled:{type:Boolean,default:void 0},indeterminate:Boolean,label:String,focusable:{type:Boolean,default:!0},checkedValue:{type:[Boolean,String,Number],default:!0},uncheckedValue:{type:[Boolean,String,Number],default:!1},"onUpdate:checked":[Function,Array],onUpdateChecked:[Function,Array],privateInsideTable:Boolean,onChange:[Function,Array]},W=n({name:`Checkbox`,props:U,setup(e){let t=s(G,null),n=C(null),{mergedClsPrefixRef:r,inlineThemeDisabled:i,mergedRtlRef:o,mergedComponentPropsRef:c}=T(e),l=C(e.defaultChecked),u=w(e,`checked`),f=j(u,l),p=N(()=>{if(t){let n=t.valueSetRef.value;return n&&e.value!==void 0?n.has(e.value):!1}return f.value===e.checkedValue}),m=P(e,{mergedSize(n){let{size:r}=e;if(r!==void 0)return r;if(t){let{value:e}=t.mergedSizeRef;if(e!==void 0)return e}if(n){let{mergedSize:e}=n;if(e!==void 0)return e.value}return c?.value?.Checkbox?.size||`medium`},mergedDisabled(n){let{disabled:r}=e;if(r!==void 0)return r;if(t){if(t.disabledRef.value)return!0;let{maxRef:{value:e},checkedCountRef:n}=t;if(e!==void 0&&n.value>=e&&!p.value)return!0;let{minRef:{value:r}}=t;if(r!==void 0&&n.value<=r&&p.value)return!0}return n?n.disabled.value:!1}}),{mergedDisabledRef:h,mergedSizeRef:g}=m,v=_(`Checkbox`,`-checkbox`,B,L,e,r);function y(n){if(t&&e.value!==void 0)t.toggleCheckbox(!p.value,e.value);else{let{onChange:t,"onUpdate:checked":r,onUpdateChecked:i}=e,{nTriggerFormInput:a,nTriggerFormChange:o}=m,s=p.value?e.uncheckedValue:e.checkedValue;r&&k(r,s,n),i&&k(i,s,n),t&&k(t,s,n),a(),o(),l.value=s}}function b(e){h.value||y(e)}function E(e){if(!h.value)switch(e.key){case` `:case`Enter`:y(e)}}function D(e){e.key===` `&&e.preventDefault()}let O={focus:()=>{n.value?.focus()},blur:()=>{n.value?.blur()}},A=x(`Checkbox`,o,r),F=a(()=>{let{value:e}=g,{common:{cubicBezierEaseInOut:t},self:{borderRadius:n,color:r,colorChecked:i,colorDisabled:a,colorTableHeader:o,colorTableHeaderModal:s,colorTableHeaderPopover:c,checkMarkColor:l,checkMarkColorDisabled:u,border:f,borderFocus:p,borderDisabled:m,borderChecked:h,boxShadowFocus:_,textColor:y,textColorDisabled:b,checkMarkColorDisabledChecked:x,colorDisabledChecked:S,borderDisabledChecked:C,labelPadding:w,labelLineHeight:T,labelFontWeight:E,[d(`fontSize`,e)]:D,[d(`size`,e)]:O}}=v.value;return{"--n-label-line-height":T,"--n-label-font-weight":E,"--n-size":O,"--n-bezier":t,"--n-border-radius":n,"--n-border":f,"--n-border-checked":h,"--n-border-focus":p,"--n-border-disabled":m,"--n-border-disabled-checked":C,"--n-box-shadow-focus":_,"--n-color":r,"--n-color-checked":i,"--n-color-table":o,"--n-color-table-modal":s,"--n-color-table-popover":c,"--n-color-disabled":a,"--n-color-disabled-checked":S,"--n-text-color":y,"--n-text-color-disabled":b,"--n-check-mark-color":l,"--n-check-mark-color-disabled":u,"--n-check-mark-color-disabled-checked":x,"--n-font-size":D,"--n-label-padding":w}}),I=i?S(`checkbox`,a(()=>g.value[0]),F,e):void 0;return Object.assign(m,O,{rtlEnabled:A,selfRef:n,mergedClsPrefix:r,mergedDisabled:h,renderedChecked:p,mergedTheme:v,labelId:M(),handleClick:b,handleKeyUp:E,handleKeyDown:D,cssVars:i?void 0:F,themeClass:I?.themeClass,onRender:I?.onRender})},render(){let{$slots:e,renderedChecked:t,mergedDisabled:n,indeterminate:a,privateInsideTable:s,cssVars:c,labelId:l,label:u,mergedClsPrefix:d,focusable:f,handleKeyUp:g,handleKeyDown:_,handleClick:v}=this;this.onRender?.();let b=A(e.default,e=>u||e?(r(),i(`span`,{key:1,class:p(`${d}-checkbox__label`),id:l},[D(()=>u||e)],10,V)):null);return(()=>{let e=m(`70be6e74cd27cb50`);return r(),i(`div`,{ref:`selfRef`,class:p([`${d}-checkbox`,this.themeClass,this.rtlEnabled&&`${d}-checkbox--rtl`,t&&`${d}-checkbox--checked`,n&&`${d}-checkbox--disabled`,a&&`${d}-checkbox--indeterminate`,s&&`${d}-checkbox--inside-table`,b&&`${d}-checkbox--show-label`]),tabindex:n||!f?void 0:0,role:`checkbox`,"aria-checked":a?`mixed`:t,"aria-labelledby":l,style:y(c),onKeyup:g,onKeydown:_,onClick:v,onMousedown:e[0]||=()=>{O(`selectstart`,window,e=>{e.preventDefault()},{once:!0})}},[h(`div`,{class:p(`${d}-checkbox-box-wrapper`)},[e[1]||=D(`\xA0`,-1),h(`div`,{class:p(`${d}-checkbox-box`)},[o(I,null,{default:()=>this.indeterminate?(r(),i(`div`,{key:`indeterminate`,class:p(`${d}-checkbox-icon`)},[D(()=>z())],2)):(r(),i(`div`,{key:`check`,class:p(`${d}-checkbox-icon`)},[D(()=>R())],2))},1024),h(`div`,{class:p(`${d}-checkbox-box__border`)},null,2)],2)],2),D(()=>b)],46,H)})()}}),G=t(`n-checkbox-group`),K=n({name:`CheckboxGroup`,props:{min:Number,max:Number,size:String,options:Array,labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},value:Array,defaultValue:{type:Array,default:null},disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]},setup(e){let{mergedClsPrefixRef:t}=T(e),n=P(e),{mergedSizeRef:r,mergedDisabledRef:i}=n,o=C(e.defaultValue),s=a(()=>e.value),l=j(s,o),u=a(()=>l.value?.length||0),d=a(()=>Array.isArray(l.value)?new Set(l.value):new Set);function f(t,r){let{nTriggerFormInput:i,nTriggerFormChange:a}=n,{onChange:s,"onUpdate:value":c,onUpdateValue:u}=e;if(Array.isArray(l.value)){let e=Array.from(l.value),n=e.findIndex(e=>e===r);t?~n||(e.push(r),u&&k(u,e,{actionType:`check`,value:r}),c&&k(c,e,{actionType:`check`,value:r}),i(),a(),o.value=e,s&&k(s,e)):~n&&(e.splice(n,1),u&&k(u,e,{actionType:`uncheck`,value:r}),c&&k(c,e,{actionType:`uncheck`,value:r}),s&&k(s,e),o.value=e,i(),a())}else t?(u&&k(u,[r],{actionType:`check`,value:r}),c&&k(c,[r],{actionType:`check`,value:r}),s&&k(s,[r]),o.value=[r],i(),a()):(u&&k(u,[],{actionType:`uncheck`,value:r}),c&&k(c,[],{actionType:`uncheck`,value:r}),s&&k(s,[]),o.value=[],i(),a())}return c(G,{checkedCountRef:u,maxRef:w(e,`max`),minRef:w(e,`min`),valueSetRef:d,disabledRef:i,mergedSizeRef:r,toggleCheckbox:f}),{mergedClsPrefix:t}},render(){let{options:t,labelField:n,valueField:a}=this.$props;return r(),i(`div`,{class:p(`${this.mergedClsPrefix}-checkbox-group`),role:`group`},[t?(r(),i(E,{key:0},[D(()=>t.map(t=>{let i=t[a];return r(),e(W,{key:i,value:i,disabled:t.disabled,label:t[n]},null,8,[`value`,`disabled`,`label`])}))],64)):(r(),i(E,{key:1},[D(()=>this.$slots.default?.())],64))],2)}});export{W as n,K as t};