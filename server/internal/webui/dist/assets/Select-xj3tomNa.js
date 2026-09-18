import{$t as e,Bt as t,D as n,E as r,Et as i,Gt as a,Ht as o,It as s,Jt as c,Kt as l,Lt as u,Mt as d,P as f,Qt as p,Rt as m,T as h,Ut as g,Xt as _,Yt as v,a as y,bt as b,cn as x,d as S,en as C,f as w,gt as T,h as ee,hn as E,ht as D,j as O,k,kt as A,l as j,m as M,nn as N,rn as P,sn as F,u as te,un as I,vn as L,vt as R,xn as z,xt as B,yt as V}from"./vue-i18n-BxOqVJSv.js";import{c as ne,n as H,o as U,r as W,s as G}from"./event-Ef8K3FrB.js";import{c as K,d as q,f as J,g as re,h as ie,l as ae,m as oe,n as se,o as Y,p as ce,r as X,s as le,t as ue,u as de}from"./create-CKNSf1FL.js";import{D as fe,F as pe,N as me,O as Z,P as Q,T as he,g as ge,j as _e,m as ve,v as ye,x as be,y as xe}from"./_plugin-vue_export-helper-3sd4LcYN.js";import{a as Se,f as Ce,t as we,x as Te}from"./fade-in-scale-up.cssr-DpVYUTqJ.js";import{C as Ee,E as $,I as De,b as Oe}from"./index-0G3VD5Ph.js";function ke(e){return e&-e}var Ae=class{constructor(e,t){this.l=e,this.min=t;let n=Array(e+1);for(let t=0;t<e+1;++t)n[t]=0;this.ft=n}add(e,t){if(t===0)return;let{l:n,ft:r}=this;for(e+=1;e<=n;)r[e]+=t,e+=ke(e)}get(e){return this.sum(e+1)-this.sum(e)}sum(e){if(e===void 0&&(e=this.l),e<=0)return 0;let{ft:t,min:n,l:r}=this;if(e>r)throw Error("[FinweckTree.sum]: `i` is larger than length.");let i=e*n;for(;e>0;)i+=t[e],e-=ke(e);return i}getBound(e){let t=0,n=this.l;for(;n>t;){let r=Math.floor((t+n)/2),i=this.sum(r);if(i>e){n=r;continue}if(i<e){if(t===r)return this.sum(t+1)<=e?t+1:r;t=r}else return r}return t}},je;function Me(){return typeof document>`u`?!1:(je===void 0&&(je=`matchMedia`in window&&window.matchMedia(`(pointer:coarse)`).matches),je)}var Ne;function Pe(){return typeof document>`u`?1:(Ne===void 0&&(Ne=`chrome`in window?window.devicePixelRatio:1),Ne)}var Fe=`VVirtualListXScroll`;function Ie({columnsRef:e,renderColRef:t,renderItemWithColsRef:n}){let r=E(0),i=E(0),a=s(()=>{let t=e.value;if(t.length===0)return null;let n=new Ae(t.length,0);return t.forEach((e,t)=>{n.add(t,e.width)}),n}),o=S(()=>{let e=a.value;return e===null?0:Math.max(e.getBound(i.value)-1,0)}),c=e=>{let t=a.value;return t===null?0:t.sum(e)},l=S(()=>{let t=a.value;return t===null?0:Math.min(t.getBound(i.value+r.value)+1,e.value.length-1)});return P(Fe,{startIndexRef:o,endIndexRef:l,columnsRef:e,renderColRef:t,renderItemWithColsRef:n,getLeft:c}),{listWidthRef:r,scrollLeftRef:i}}var Le=g({name:`VirtualListRow`,props:{index:{type:Number,required:!0},item:{type:Object,required:!0}},setup(){let{startIndexRef:e,endIndexRef:t,columnsRef:n,getLeft:r,renderColRef:i,renderItemWithColsRef:a}=l(Fe);return{startIndex:e,endIndex:t,columns:n,renderCol:i,renderItemWithCols:a,getLeft:r}},render(){let{startIndex:e,endIndex:t,columns:n,renderCol:r,renderItemWithCols:i,getLeft:a,item:o}=this;if(i!=null)return i({itemIndex:this.index,startColIndex:e,endColIndex:t,allColumns:n,item:o,getLeft:a});if(r!=null){let i=[];for(let s=e;s<=t;++s){let e=n[s];i.push(r({column:e,left:a(s),item:o}))}return i}return null}}),Re=K(`.v-vl`,{maxHeight:`inherit`,height:`100%`,overflow:`auto`,minWidth:`1px`},[K(`&:not(.v-vl--show-scrollbar)`,{scrollbarWidth:`none`},[K(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,{width:0,height:0,display:`none`})])]),ze=g({name:`VirtualList`,inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},columns:{type:Array,default:()=>[]},renderCol:Function,renderItemWithCols:Function,items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:`div`},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:`key`},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(t){let n=O();Re.mount({id:`vueuc/virtual-list`,head:!0,anchorMetaName:ae,ssr:n}),C(()=>{let{defaultScrollIndex:e,defaultScrollKey:n}=t;e==null?n!=null&&y({key:n}):y({index:e})});let r=!1,i=!1;_(()=>{if(r=!1,!i){i=!0;return}y({top:h.value,left:c.value})}),e(()=>{r=!0,i||=!0});let a=S(()=>{if(t.renderCol==null&&t.renderItemWithCols==null||t.columns.length===0)return;let e=0;return t.columns.forEach(t=>{e+=t.width}),e}),o=s(()=>{let e=new Map,{keyField:n}=t;return t.items.forEach((t,r)=>{e.set(t[n],r)}),e}),{scrollLeftRef:c,listWidthRef:l}=Ie({columnsRef:L(t,`columns`),renderColRef:L(t,`renderCol`),renderItemWithColsRef:L(t,`renderItemWithCols`)}),u=E(null),d=E(void 0),f=new Map,p=s(()=>{let{items:e,itemSize:n,keyField:r}=t,i=new Ae(e.length,n);return e.forEach((e,t)=>{let n=e[r],a=f.get(n);a!==void 0&&i.add(t,a)}),i}),m=E(0),h=E(0),g=S(()=>Math.max(p.value.getBound(h.value-_e(t.paddingTop))-1,0)),v=s(()=>{let{value:e}=d;if(e===void 0)return[];let{items:n,itemSize:r}=t,i=g.value,a=Math.min(i+Math.ceil(e/r+1),n.length-1),o=[];for(let e=i;e<=a;++e)o.push(n[e]);return o}),y=(e,t)=>{if(typeof e==`number`){T(e,t,`auto`);return}let{left:n,top:r,index:i,key:a,position:s,behavior:c,debounce:l=!0}=e;if(n!==void 0||r!==void 0)T(n,r,c);else if(i!==void 0)w(i,c,l);else if(a!==void 0){let e=o.value.get(a);e!==void 0&&w(e,c,l)}else s===`bottom`?T(0,2**53-1,c):s===`top`&&T(0,0,c)},b,x=null;function w(e,n,r){let i=u.value;if(i==null)return;let{value:a}=p,o=a.sum(e)+_e(t.paddingTop);if(!r)i.scrollTo({left:0,top:o,behavior:n});else{b=e,x!==null&&window.clearTimeout(x),x=window.setTimeout(()=>{b=void 0,x=null},16);let{scrollTop:t,offsetHeight:r}=i;if(o>t){let s=a.get(e);o+s<=t+r||i.scrollTo({left:0,top:o+s-r,behavior:n})}else i.scrollTo({left:0,top:o,behavior:n})}}function T(e,t,n){u.value?.scrollTo({left:e,top:t,behavior:n})}function ee(e,n){if(r||t.ignoreItemResize||P(n.target))return;let{value:i}=p,a=o.value.get(e),s=i.get(a),c=n.borderBoxSize?.[0]?.blockSize??n.contentRect.height;if(c===s)return;c-t.itemSize===0?f.delete(e):f.set(e,c-t.itemSize);let l=c-s;if(l===0)return;i.add(a,l);let d=u.value;if(d!=null){if(b===void 0){let e=i.sum(a);d.scrollTop>e&&d.scrollBy(0,l)}else(a<b||a===b&&c+i.sum(a)>d.scrollTop+d.offsetHeight)&&d.scrollBy(0,l);N()}m.value++}let D=!Me(),k=!1;function A(e){var n;(n=t.onScroll)==null||n.call(t,e),(!D||!k)&&N()}function j(e){var n;if((n=t.onWheel)==null||n.call(t,e),D){let t=u.value;if(t!=null){if(e.deltaX===0&&(t.scrollTop===0&&e.deltaY<=0||t.scrollTop+t.offsetHeight>=t.scrollHeight&&e.deltaY>=0))return;e.preventDefault(),t.scrollTop+=e.deltaY/Pe(),t.scrollLeft+=e.deltaX/Pe(),N(),k=!0,Te(()=>{k=!1})}}}function M(e){if(r||P(e.target))return;if(t.renderCol==null&&t.renderItemWithCols==null){if(e.contentRect.height===d.value)return}else if(e.contentRect.height===d.value&&e.contentRect.width===l.value)return;d.value=e.contentRect.height,l.value=e.contentRect.width;let{onResize:n}=t;n!==void 0&&n(e)}function N(){let{value:e}=u;e!=null&&(h.value=e.scrollTop,c.value=e.scrollLeft)}function P(e){let t=e;for(;t!==null;){if(t.style.display===`none`)return!0;t=t.parentElement}return!1}return{listHeight:d,listStyle:{overflow:`auto`},keyToIndex:o,itemsStyle:s(()=>{let{itemResizable:e}=t,n=Q(p.value.sum());return m.value,[t.itemsStyle,{boxSizing:`content-box`,width:Q(a.value),height:e?``:n,minHeight:e?n:``,paddingTop:Q(t.paddingTop),paddingBottom:Q(t.paddingBottom)}]}),visibleItemsStyle:s(()=>(m.value,{transform:`translateY(${Q(p.value.sum(g.value))})`})),viewportItems:v,listElRef:u,itemsElRef:E(null),scrollTo:y,handleListResize:M,handleListScroll:A,handleListWheel:j,handleItemResize:ee}},render(){let{itemResizable:e,keyField:t,keyToIndex:n,visibleItemsTag:r}=this;return a(ye,{onResize:this.handleListResize},{default:()=>{var i;return a(`div`,c(this.$attrs,{class:[`v-vl`,this.showScrollbar&&`v-vl--show-scrollbar`],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:`listElRef`}),[this.items.length===0?(i=this.$slots).empty?.call(i):a(`div`,{ref:`itemsElRef`,class:`v-vl-items`,style:this.itemsStyle},[a(r,Object.assign({class:`v-vl-visible-items`,style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>{let{renderCol:r,renderItemWithCols:i}=this;return this.viewportItems.map(o=>{let s=o[t],c=n.get(s),l=r==null?void 0:a(Le,{index:c,item:o}),u=i==null?void 0:a(Le,{index:c,item:o}),d=this.$slots.default({item:o,renderedCols:l,renderedItemWithCols:u,index:c})[0];return e?a(ye,{key:s,onResize:e=>this.handleItemResize(s,e)},{default:()=>d}):(d.key=s,d)})}})])])}})}});function Be(e){switch(typeof e){case`string`:return e||void 0;case`number`:return String(e);default:return}}function Ve(e,t){t&&(C(()=>{let{value:n}=e;n&&xe.registerHandler(n,t)}),F(e,(e,t)=>{t&&xe.unregisterHandler(t)},{deep:!1}),p(()=>{let{value:t}=e;t&&xe.unregisterHandler(t)}))}var He=g({props:{onFocus:Function,onBlur:Function},setup(e){return()=>(()=>{let n=h(`d16ead82505dc285`);return N(),t(`div`,{style:`width: 0; height: 0`,tabindex:0,onFocus:n[0]||=t=>e.onFocus?.(t),onBlur:n[1]||=t=>e.onBlur?.(t)},null,32)})()}}),Ue=g({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=l(ie);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:n,renderOption:r,nodeProps:i,tmNode:{rawNode:a}}=this,o=i?.(a),s=n?n(a,!1):$(a[this.labelField],a,!1),l=(N(),t(`div`,c(o,{class:[`${e}-base-select-group-header`,o?.class]}),[k(()=>s)],16));return a.render?a.render({node:l,option:a}):r?r({node:l,option:a,selected:!1}):l}});function We(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var Ge=g({name:`Checkmark`,render(){return(()=>{let e=h(`3c84eac8ae4e1f96`);return e[0]||=u(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},[u(`g`,{fill:`none`},[u(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})])],-1)})()}}),Ke=[`onClick`,`onMouseenter`,`onMousemove`];function qe(e,t){return N(),m(i,{name:`fade-in-scale-up-transition`},{default:()=>e?(N(),m(w,{key:1,clsPrefix:t,class:r(`${t}-base-select-option__check`)},{default:()=>a(Ge)},1032,[`clsPrefix`,`class`])):null},1024)}var Je=g({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:u,nodePropsRef:d,handleOptionClick:f,handleOptionMouseEnter:p}=l(ie),m=S(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e;n.disabled||p(t,n)}function _(t){let{tmNode:n}=e,{value:r}=m;n.disabled||r||p(t,n)}return{multiple:r,isGrouped:S(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:u,nodeProps:d,isPending:m,isSelected:S(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:_,handleMouseEnter:g,handleClick:h}},render(){let{clsPrefix:e,tmNode:{rawNode:n},isSelected:i,isPending:a,isGrouped:o,showCheckmark:s,nodeProps:l,renderOption:d,renderLabel:f,handleClick:p,handleMouseEnter:m,handleMouseMove:h}=this,g=qe(i,e),_=f?[f(n,i),s&&g]:[$(n[this.labelField],n,i),s&&g],v=l?.(n),y=(N(),t(`div`,c(v,{class:[`${e}-base-select-option`,n.class,v?.class,{[`${e}-base-select-option--disabled`]:n.disabled,[`${e}-base-select-option--selected`]:i,[`${e}-base-select-option--grouped`]:o,[`${e}-base-select-option--pending`]:a,[`${e}-base-select-option--show-checkmark`]:s}],style:[v?.style||``,n.style||``],onClick:We([p,v?.onClick]),onMouseenter:We([m,v?.onMouseenter]),onMousemove:We([h,v?.onMousemove])}),[u(`div`,{class:r(`${e}-base-select-option__content`)},[k(()=>_)],2)],16,Ke));return n.render?n.render({node:y,option:n,selected:i}):d?d({node:y,option:n,selected:i}):y}}),Ye=T(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[T(`scrollbar`,`
 max-height: var(--n-height);
 `),T(`virtual-list`,`
 max-height: var(--n-height);
 `),T(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[R(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),T(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),T(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),R(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),R(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),R(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),R(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),T(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),T(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[V(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),D(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),D(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),V(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),V(`pending`,[D(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),V(`selected`,`
 color: var(--n-option-text-color-active);
 `,[D(`&::before`,`
 background-color: var(--n-option-color-active);
 `),V(`pending`,[D(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),V(`disabled`,`
 cursor: not-allowed;
 `,[b(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),V(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),R(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[we({enterScale:`0.5`})])])]),Xe=[`tabindex`,`onFocusin`,`onFocusout`,`onKeyup`,`onKeydown`,`onMousedown`,`onMouseenter`,`onMouseleave`],Ze=g({name:`InternalSelectMenu`,props:{...M.props,clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function},setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n,mergedComponentPropsRef:r}=f(e),i=j(`InternalSelectMenu`,n,t),a=M(`InternalSelectMenu`,`-internal-select-menu`,Ye,De,e,L(e,`clsPrefix`)),o=E(null),c=E(null),l=E(null),u=s(()=>e.treeMate.getFlattenedNodes()),d=s(()=>se(u.value)),m=E(null);function h(){let{treeMate:t}=e,n=null,{value:r}=e;r===null?n=t.getFirstAvailableNode():(n=e.multiple?t.getNode((r||[])[(r||[]).length-1]):t.getNode(r),(!n||n.disabled)&&(n=t.getFirstAvailableNode())),H(n||null)}function g(){let{value:t}=m;t&&!e.treeMate.getNode(t.key)&&(m.value=null)}let _;F(()=>e.show,t=>{t?_=F(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?h():g(),v(U)):g()},{immediate:!0}):_?.()},{immediate:!0}),p(()=>{_?.()});let y=s(()=>_e(a.value.self[B(`optionHeight`,e.size)])),b=s(()=>me(a.value.self[B(`padding`,e.size)])),x=s(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),S=s(()=>{let e=u.value;return e&&e.length===0}),w=s(()=>r?.value?.Select?.renderEmpty);function T(t){let{onToggle:n}=e;n&&n(t)}function D(t){let{onScroll:n}=e;n&&n(t)}function O(e){l.value?.sync(),D(e)}function k(){l.value?.sync()}function A(){let{value:e}=m;return e||null}function N(e,t){t.disabled||H(t,!1)}function te(e,t){t.disabled||T(t)}function I(t){re(t,`action`)||e.onKeyup?.(t)}function R(t){re(t,`action`)||e.onKeydown?.(t)}function z(t){e.onMousedown?.(t),!e.focusable&&t.preventDefault()}function V(){let{value:e}=m;e&&H(e.getNext({loop:!0}),!0)}function ne(){let{value:e}=m;e&&H(e.getPrev({loop:!0}),!0)}function H(e,t=!1){m.value=e,t&&U()}function U(){let t=m.value;if(!t)return;let n=d.value(t.key);n!==null&&(e.virtualScroll?c.value?.scrollTo({index:n}):l.value?.scrollTo({index:n,elSize:y.value}))}function W(t){o.value?.contains(t.target)&&e.onFocus?.(t)}function G(t){o.value?.contains(t.relatedTarget)||e.onBlur?.(t)}P(ie,{handleOptionMouseEnter:N,handleOptionClick:te,valueSetRef:x,pendingTmNodeRef:m,nodePropsRef:L(e,`nodeProps`),showCheckmarkRef:L(e,`showCheckmark`),multipleRef:L(e,`multiple`),valueRef:L(e,`value`),renderLabelRef:L(e,`renderLabel`),renderOptionRef:L(e,`renderOption`),labelFieldRef:L(e,`labelField`),valueFieldRef:L(e,`valueField`)}),P(oe,o),C(()=>{let{value:e}=l;e&&e.sync()});let K=s(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:o,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[B(`optionFontSize`,t)]:x,[B(`optionHeight`,t)]:S,[B(`optionPadding`,t)]:C}}=a.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":o,"--n-option-font-size":x,"--n-group-header-text-color":s,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":me(C,`left`),"--n-option-padding-right":me(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:q}=e,J=q?ee(`internal-select-menu`,s(()=>e.size[0]),K,e):void 0,ae={selfRef:o,next:V,prev:ne,getPendingTmNode:A};return Ve(o,e.onResize),{mergedTheme:a,mergedClsPrefix:t,rtlEnabled:i,virtualListRef:c,scrollbarRef:l,itemSize:y,padding:b,flattenedNodes:u,empty:S,mergedRenderEmpty:w,virtualListContainer(){let{value:e}=c;return e?.listElRef},virtualListContent(){let{value:e}=c;return e?.itemsElRef},doScroll:D,handleFocusin:W,handleFocusout:G,handleKeyUp:I,handleKeyDown:R,handleMouseDown:z,handleVirtualListResize:k,handleVirtualListScroll:O,cssVars:q?void 0:K,themeClass:J?.themeClass,onRender:J?.onRender,...ae}},render(){let{$slots:e,virtualScroll:n,clsPrefix:i,mergedTheme:a,themeClass:o,onRender:s}=this;return s?.(),N(),t(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:r([`${i}-base-select-menu`,`${i}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${i}-base-select-menu--rtl`,o,this.multiple&&`${i}-base-select-menu--multiple`]),style:z(this.cssVars),onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},[k(()=>fe(e.header,e=>e&&(N(),t(`div`,{class:r(`${i}-base-select-menu__header`),"data-header":!0,key:`header`},[k(()=>e)],2)))),this.loading?(N(),t(`div`,{key:0,class:r(`${i}-base-select-menu__loading`)},[(N(),m(y,{clsPrefix:i,strokeWidth:20},null,8,[`clsPrefix`]))],2)):(N(),t(d,{key:1},[this.empty?(N(),t(`div`,{key:1,class:r(`${i}-base-select-menu__empty`),"data-empty":!0},[k(()=>he(e.empty,()=>[this.mergedRenderEmpty?.()||(N(),m(G,{theme:a.peers.Empty,themeOverrides:a.peerOverrides.Empty,size:this.size},null,8,[`theme`,`themeOverrides`,`size`]))]))],2)):(N(),m(ge,c({key:0,ref:`scrollbarRef`,theme:a.peers.Scrollbar,themeOverrides:a.peerOverrides.Scrollbar,scrollable:this.scrollable,container:n?this.virtualListContainer:void 0,content:n?this.virtualListContent:void 0,onScroll:n?void 0:this.doScroll},this.scrollbarProps),{default:()=>n?(N(),m(ze,{key:1,ref:`virtualListRef`,class:r(`${i}-virtual-list`),items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?(N(),m(Ue,{key:e.key,clsPrefix:i,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):e.ignored?null:(N(),m(Je,{clsPrefix:i,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))},1032,[`class`,`items`,`itemSize`,`paddingTop`,`paddingBottom`,`onResize`,`onScroll`])):(N(),t(`div`,{key:4,class:r(`${i}-base-select-menu-option-wrapper`),style:z({paddingTop:this.padding.top,paddingBottom:this.padding.bottom})},[k(()=>this.flattenedNodes.map(e=>e.isGroup?(N(),m(Ue,{key:e.key,clsPrefix:i,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):(N(),m(Je,{clsPrefix:i,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))))],6))},1040,[`theme`,`themeOverrides`,`scrollable`,`container`,`content`,`onScroll`]))],64)),k(()=>fe(e.action,e=>e&&[(N(),t(`div`,{class:r(`${i}-base-select-menu__action`),"data-action":!0,key:`action`},[k(()=>e)],2)),(N(),m(He,{onFocus:this.onTabOut,key:`focus-detector`},null,8,[`onFocus`]))]))],46,Xe)}});function Qe(e){return e.type===`group`}function $e(e){return e.type===`ignored`}function et(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function tt(e,t){return{getIsGroup:Qe,getIgnored:$e,getKey(t){return Qe(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function nt(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(Qe(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if($e(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function rt(e,t,n){let r=new Map;return e.forEach(e=>{Qe(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var it=D([T(`base-selection`,`
 --n-padding-single: var(--n-padding-single-top) var(--n-padding-single-right) var(--n-padding-single-bottom) var(--n-padding-single-left);
 --n-padding-multiple: var(--n-padding-multiple-top) var(--n-padding-multiple-right) var(--n-padding-multiple-bottom) var(--n-padding-multiple-left);
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[T(`base-loading`,`
 color: var(--n-loading-color);
 `),T(`base-selection-tags`,`min-height: var(--n-height);`),R(`border, state-border`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),R(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),T(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[R(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),T(`base-selection-overlay`,`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[R(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),T(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[R(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),T(`base-selection-tags`,`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),T(`base-selection-label`,`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[T(`base-selection-input`,`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[R(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),R(`render-label`,`
 color: var(--n-text-color);
 `)]),b(`disabled`,[D(`&:hover`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),V(`focus`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),V(`active`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),T(`base-selection-label`,`background-color: var(--n-color-active);`),T(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),V(`disabled`,`cursor: not-allowed;`,[R(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),T(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[T(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),R(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),T(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),T(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),T(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[R(`input`,`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),R(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>V(`${e}-status`,[R(`state-border`,`border: var(--n-border-${e});`),b(`disabled`,[D(`&:hover`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),V(`active`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),T(`base-selection-label`,`background-color: var(--n-color-active-${e});`),T(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),V(`focus`,[R(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),T(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),T(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[D(`&:last-child`,`padding-right: 0;`),T(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[R(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),at=[`disabled`,`value`,`autofocus`,`onBlur`,`onFocus`,`onKeydown`,`onInput`,`onCompositionstart`,`onCompositionend`],ot=[`tabindex`],st=[`title`],ct=[`value`,`readonly`,`disabled`,`autofocus`,`onFocus`,`onBlur`,`onInput`,`onCompositionstart`,`onCompositionend`],lt=[`tabindex`],ut=[`onClick`,`onMouseenter`,`onMouseleave`,`onKeydown`,`onFocusin`,`onFocusout`,`onMousedown`],dt=g({name:`InternalSelection`,props:{...M.props,clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function},setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n}=f(e),r=j(`InternalSelection`,n,t),i=E(null),a=E(null),o=E(null),c=E(null),l=E(null),u=E(null),d=E(null),p=E(null),m=E(null),h=E(null),g=E(!1),_=E(!1),y=E(!1),b=M(`InternalSelection`,`-internal-selection`,it,Ee,e,L(e,`clsPrefix`)),S=s(()=>e.clearable&&!e.disabled&&(y.value||e.active)),w=s(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):$(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),T=s(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),D=s(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function O(){let{value:t}=i;if(t){let{value:n}=a;n&&(n.style.width=`${t.offsetWidth}px`,e.maxTagCount!==`responsive`&&m.value?.sync({showAllItemsBeforeCalculate:!1}))}}function k(){let{value:e}=h;e&&(e.style.display=`none`)}function A(){let{value:e}=h;e&&(e.style.display=`inline-block`)}F(L(e,`active`),e=>{e||k()}),F(L(e,`pattern`),()=>{e.multiple&&v(O)});function N(t){let{onFocus:n}=e;n&&n(t)}function P(t){let{onBlur:n}=e;n&&n(t)}function te(t){let{onDeleteOption:n}=e;n&&n(t)}function I(t){let{onClear:n}=e;n&&n(t)}function R(t){let{onPatternInput:n}=e;n&&n(t)}function z(e){(!e.relatedTarget||!o.value?.contains(e.relatedTarget))&&N(e)}function V(e){o.value?.contains(e.relatedTarget)||P(e)}function ne(e){I(e)}function H(){y.value=!0}function U(){y.value=!1}function W(t){e.active&&e.filterable&&t.target!==a.value&&t.preventDefault()}function G(e){te(e)}let K=E(!1);function q(t){if(t.key===`Backspace`&&!K.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&G(t[t.length-1])}}let J=null;function re(t){let{value:n}=i;n&&(n.textContent=t.target.value,O()),e.ignoreComposition&&K.value?J=t:R(t)}function ie(){K.value=!0}function ae(){K.value=!1,e.ignoreComposition&&R(J),J=null}function oe(t){_.value=!0,e.onPatternFocus?.(t)}function se(t){_.value=!1,e.onPatternBlur?.(t)}function Y(){if(e.filterable)_.value=!1,u.value?.blur(),a.value?.blur();else if(e.multiple){let{value:e}=c;e?.blur()}else{let{value:e}=l;e?.blur()}}function ce(){e.filterable?(_.value=!1,u.value?.focus()):e.multiple?c.value?.focus():l.value?.focus()}function X(){let{value:e}=a;e&&(A(),e.focus())}function le(){let{value:e}=a;e&&e.blur()}function ue(e){let{value:t}=d;t&&t.setTextContent(`+${e}`)}function de(){let{value:e}=p;return e}function fe(){return a.value}let pe=null;function Z(){pe!==null&&window.clearTimeout(pe)}function Q(){e.active||(Z(),pe=window.setTimeout(()=>{D.value&&(g.value=!0)},100))}function he(){Z()}function ge(e){e||(Z(),g.value=!1)}F(D,e=>{e||(g.value=!1)}),C(()=>{x(()=>{let t=u.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=_.value?-1:0)})}),Ve(o,e.onResize);let{inlineThemeDisabled:_e}=e,ve=s(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:x,borderActive:S,arrowColor:C,arrowColorDisabled:w,loadingColor:T,colorActiveWarning:ee,boxShadowFocusWarning:E,boxShadowActiveWarning:D,boxShadowHoverWarning:O,borderWarning:k,borderFocusWarning:A,borderHoverWarning:j,borderActiveWarning:M,colorActiveError:N,boxShadowFocusError:P,boxShadowActiveError:F,boxShadowHoverError:te,borderError:I,borderFocusError:L,borderHoverError:R,borderActiveError:z,clearColor:V,clearColorHover:ne,clearColorPressed:H,clearSize:U,arrowSize:W,[B(`height`,t)]:G,[B(`fontSize`,t)]:K}}=b.value,q=me(c),J=me(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":y,"--n-border-hover":x,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":K,"--n-height":G,"--n-padding-single-top":q.top,"--n-padding-multiple-top":J.top,"--n-padding-single-right":q.right,"--n-padding-multiple-right":J.right,"--n-padding-single-left":q.left,"--n-padding-multiple-left":J.left,"--n-padding-single-bottom":q.bottom,"--n-padding-multiple-bottom":J.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":C,"--n-arrow-color-disabled":w,"--n-loading-color":T,"--n-color-active-warning":ee,"--n-box-shadow-focus-warning":E,"--n-box-shadow-active-warning":D,"--n-box-shadow-hover-warning":O,"--n-border-warning":k,"--n-border-focus-warning":A,"--n-border-hover-warning":j,"--n-border-active-warning":M,"--n-color-active-error":N,"--n-box-shadow-focus-error":P,"--n-box-shadow-active-error":F,"--n-box-shadow-hover-error":te,"--n-border-error":I,"--n-border-focus-error":L,"--n-border-hover-error":R,"--n-border-active-error":z,"--n-clear-size":U,"--n-clear-color":V,"--n-clear-color-hover":ne,"--n-clear-color-pressed":H,"--n-arrow-size":W,"--n-font-weight":r}}),ye=_e?ee(`internal-selection`,s(()=>e.size[0]),ve,e):void 0;return{mergedTheme:b,mergedClearable:S,mergedClsPrefix:t,rtlEnabled:r,patternInputFocused:_,filterablePlaceholder:w,label:T,selected:D,showTagsPanel:g,isComposing:K,counterRef:d,counterWrapperRef:p,patternInputMirrorRef:i,patternInputRef:a,selfRef:o,multipleElRef:c,singleElRef:l,patternInputWrapperRef:u,overflowRef:m,inputTagElRef:h,handleMouseDown:W,handleFocusin:z,handleClear:ne,handleMouseEnter:H,handleMouseLeave:U,handleDeleteOption:G,handlePatternKeyDown:q,handlePatternInputInput:re,handlePatternInputBlur:se,handlePatternInputFocus:oe,handleMouseEnterCounter:Q,handleMouseLeaveCounter:he,handleFocusout:V,handleCompositionEnd:ae,handleCompositionStart:ie,onPopoverUpdateShow:ge,focus:ce,focusInput:X,blur:Y,blurInput:le,updateCounter:ue,getCounter:de,getTail:fe,renderLabel:e.renderLabel,cssVars:_e?void 0:ve,themeClass:ye?.themeClass,onRender:ye?.onRender}},render(){let{status:e,multiple:n,size:i,disabled:a,filterable:o,maxTagCount:s,bordered:l,clsPrefix:f,ellipsisTagPopoverProps:p,onRender:h,renderTag:g,renderLabel:_}=this;h?.();let v=s===`responsive`,y=typeof s==`number`,b=v||y,x=(N(),m(be,null,{default:()=>(N(),m(W,{clsPrefix:f,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>this.$slots.arrow?.()},1032,[`clsPrefix`,`loading`,`showArrow`,`showClear`,`onClear`]))},1024)),S;if(n){let{labelField:e}=this,n=n=>(N(),t(`div`,{class:r(`${f}-base-selection-tag-wrapper`),key:n.value},[g?(N(),t(d,{key:0},[k(()=>g({option:n,handleClose:()=>{this.handleDeleteOption(n)}}))],64)):(N(),m(U,{key:1,size:i,closable:!n.disabled,disabled:a,onClose:()=>{this.handleDeleteOption(n)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>_?_(n,!0):$(n[e],n,!0)},1032,[`size`,`closable`,`disabled`,`onClose`]))],2)),l=()=>(y?this.selectedOptions.slice(0,s):this.selectedOptions).map(n),h=o?(N(),t(`div`,{class:r(`${f}-base-selection-input-tag`),ref:`inputTagElRef`,key:`__input-tag__`},[u(`input`,c(this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:a,value:this.pattern,autofocus:this.autofocus,class:`${f}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,at),u(`span`,{ref:`patternInputMirrorRef`,class:r(`${f}-base-selection-input-tag__mirror`)},[k(()=>this.pattern)],2)],2)):null,C=v?()=>(N(),t(`div`,{class:r(`${f}-base-selection-tag-wrapper`),ref:`counterWrapperRef`},[(N(),m(U,{size:i,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:a},null,8,[`size`,`onMouseenter`,`onMouseleave`,`disabled`]))],2)):void 0,w;if(y){let e=this.selectedOptions.length-s;e>0&&(w=(n=>(N(),t(`div`,{class:r(`${f}-base-selection-tag-wrapper`),key:`__counter__`},[(N(),m(U,{size:i,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:a},{default:()=>`+${e}`},1032,[`size`,`onMouseenter`,`disabled`]))],2)))(w))}let T=v?o?(N(),m(Y,{key:3,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:l,counter:C,tail:()=>h},1032,[`updateCounter`,`getCounter`,`getTail`])):(N(),m(Y,{key:4,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:l,counter:C},1032,[`updateCounter`,`getCounter`])):y&&w?l().concat(w):l(),ee=b?()=>(N(),t(`div`,{class:r(`${f}-base-selection-popover`)},[v?(N(),t(d,{key:0},[k(()=>l())],64)):(N(),t(d,{key:1},[k(()=>this.selectedOptions.map(n))],64))],2)):void 0,E=b?{show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover,...p}:null,D=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?(N(),t(`div`,{key:5,class:r(`${f}-base-selection-placeholder ${f}-base-selection-overlay`)},[u(`div`,{class:r(`${f}-base-selection-placeholder__inner`)},[k(()=>this.placeholder)],2)],2)):null,O=o?(N(),t(`div`,{key:6,ref:`patternInputWrapperRef`,class:r(`${f}-base-selection-tags`)},[k(()=>T),v?k(()=>null):(N(),t(d,{key:1},[k(()=>h)],64)),k(()=>x)],2)):(N(),t(`div`,{key:7,ref:`multipleElRef`,class:r(`${f}-base-selection-tags`),tabindex:a?void 0:0},[k(()=>T),k(()=>x)],10,ot));S=(e=>(N(),t(d,{key:8},[b?(N(),m(X,c({key:0},E,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>O,default:ee},1040)):(N(),t(d,{key:1},[k(()=>O)],64)),k(()=>D)],64)))(S)}else if(o){let e=this.pattern||this.isComposing,n=this.active?!e:!this.selected,i=!this.active&&this.selected;S=(e=>(N(),t(`div`,{key:9,ref:`patternInputWrapperRef`,class:r(`${f}-base-selection-label`),title:this.patternInputFocused?void 0:Be(this.label)},[u(`input`,c(this.inputProps,{ref:`patternInputRef`,class:`${f}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:a,disabled:a,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,ct),i?(N(),t(`div`,{class:r(`${f}-base-selection-label__render-label ${f}-base-selection-overlay`),key:`input`},[u(`div`,{class:r(`${f}-base-selection-overlay__wrapper`)},[g?(N(),t(d,{key:0},[k(()=>g({option:this.selectedOption,handleClose:()=>{}}))],64)):(N(),t(d,{key:1},[_?(N(),t(d,{key:0},[k(()=>_(this.selectedOption,!0))],64)):(N(),t(d,{key:1},[k(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],2)):k(()=>null),n?(N(),t(`div`,{class:r(`${f}-base-selection-placeholder ${f}-base-selection-overlay`),key:`placeholder`},[u(`div`,{class:r(`${f}-base-selection-overlay__wrapper`)},[k(()=>this.filterablePlaceholder)],2)],2)):k(()=>null),k(()=>x)],10,st)))(S)}else S=(e=>(N(),t(`div`,{key:10,ref:`singleElRef`,class:r(`${f}-base-selection-label`),tabindex:this.disabled?void 0:0},[this.label===void 0?(N(),t(`div`,{class:r(`${f}-base-selection-placeholder ${f}-base-selection-overlay`),key:`placeholder`},[u(`div`,{class:r(`${f}-base-selection-placeholder__inner`)},[k(()=>this.placeholder)],2)],2)):(N(),t(`div`,{class:r(`${f}-base-selection-input`),title:Be(this.label),key:`input`},[u(`div`,{class:r(`${f}-base-selection-input__content`)},[g?(N(),t(d,{key:0},[k(()=>g({option:this.selectedOption,handleClose:()=>{}}))],64)):(N(),t(d,{key:1},[_?(N(),t(d,{key:0},[k(()=>_(this.selectedOption,!0))],64)):(N(),t(d,{key:1},[k(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],10,[`title`])),k(()=>x)],10,lt)))(S);return N(),t(`div`,{ref:`selfRef`,class:r([`${f}-base-selection`,this.rtlEnabled&&`${f}-base-selection--rtl`,this.themeClass,e&&`${f}-base-selection--${e}-status`,{[`${f}-base-selection--active`]:this.active,[`${f}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${f}-base-selection--disabled`]:this.disabled,[`${f}-base-selection--multiple`]:this.multiple,[`${f}-base-selection--focus`]:this.focused}]),style:z(this.cssVars),onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},[k(()=>S),l?(N(),t(`div`,{key:0,class:r(`${f}-base-selection__border`)},null,2)):k(()=>null),l?(N(),t(`div`,{key:2,class:r(`${f}-base-selection__state-border`)},null,2)):k(()=>null)],46,ut)}}),ft=D([T(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),T(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[we({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),pt={...M.props,to:J.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array},mt=g({name:`Select`,props:pt,slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:a}=f(e),o=M(`Select`,`-select`,ft,Oe,e,t),c=E(e.defaultValue),l=L(e,`value`),u=Ce(l,c),d=E(!1),p=E(``),m=ce(e,[`items`,`options`]),h=E([]),g=E([]),_=s(()=>g.value.concat(h.value).concat(m.value)),v=s(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return et(e,i);let a=t[r];return typeof a==`string`?et(e,a):typeof a==`number`&&et(e,String(a))}}),y=s(()=>{if(e.remote)return m.value;{let{value:t}=_,{value:n}=p;return!n.length||!e.filterable?t:nt(t,v.value,n,e.childrenField)}}),b=s(()=>{let{valueField:t,childrenField:n}=e,r=tt(t,n);return ue(y.value,r)}),x=s(()=>rt(_.value,e.valueField,e.childrenField)),S=E(!1),C=Ce(L(e,`show`),S),w=E(null),T=E(null),D=E(null),{localeRef:O}=ne(`Select`),k=s(()=>e.placeholder??O.value.placeholder),A=[],j=E(new Map),N=s(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function P(t){let n=e.remote,{value:r}=j,{value:i}=x,{value:a}=N,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let I=s(()=>{if(e.multiple){let{value:e}=u;return Array.isArray(e)?P(e):[]}return null}),R=s(()=>{let{value:t}=u;return!e.multiple&&!Array.isArray(t)?t===null?null:P([t])[0]||null:null}),z=ve(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:a?.value?.Select?.size||`medium`}}),{mergedSizeRef:B,mergedDisabledRef:V,mergedStatusRef:U}=z;function W(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=z;r&&Z(r,t,n),a&&Z(a,t,n),i&&Z(i,t,n),c.value=t,o(),s()}function G(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=z;n&&Z(n,t),r()}function K(){let{onClear:t}=e;t&&Z(t)}function q(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=z;n&&Z(n,t),i(),r&&Y()}function ie(t){let{onSearch:n}=e;n&&Z(n,t)}function ae(t){let{onScroll:n}=e;n&&Z(n,t)}function oe(){let{remote:t,multiple:n}=e;if(t){let{value:t}=j;if(n){let{valueField:n}=e;I.value?.forEach(e=>{t.set(e[n],e)})}else{let n=R.value;n&&t.set(n[e.valueField],n)}}}function se(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&Z(n,t),r&&Z(r,t),S.value=t}function Y(){V.value||(se(!0),S.value=!0,e.filterable&&Ne())}function X(){se(!1)}function le(){p.value=``,g.value=A}let de=E(!1);function fe(){e.filterable&&(de.value=!0)}function me(){e.filterable&&(de.value=!1,C.value||le())}function Q(){V.value||(C.value?e.filterable?Ne():X():Y())}function he(e){D.value?.selfRef?.contains(e.relatedTarget)||(d.value=!1,G(e),X())}function ge(e){q(e),d.value=!0}function _e(){d.value=!0}function ye(e){w.value?.$el.contains(e.relatedTarget)||(d.value=!1,G(e),X())}function be(){w.value?.focus(),X()}function xe(e){C.value&&(w.value?.$el.contains(pe(e))||X())}function Se(t){if(!Array.isArray(t))return[];if(N.value)return Array.from(t);{let{remote:n}=e,{value:r}=x;if(n){let{value:e}=j;return t.filter(t=>r.has(t)||e.has(t))}return t.filter(e=>r.has(e))}}function we(e){Te(e.rawNode)}function Te(t){if(V.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=g,t=e[0]||null;if(t){let e=h.value;e.length?e.push(t):h.value=[t],g.value=A}}if(r&&j.value.set(t[a],t),e.multiple){let e=Se(u.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Ee(t[a]);~e&&(h.value.splice(e,1),i&&(p.value=``))}}else e.push(t[a]),i&&(p.value=``);W(e,P(e))}else{if(n&&!r){let e=Ee(t[a]);~e?h.value=[h.value[e]]:h.value=A}Me(),X(),W(t[a],t)}}function Ee(t){return h.value.findIndex(n=>n[e.valueField]===t)}function $(t){C.value||Y();let{value:n}=t.target;p.value=n;let{tag:r,remote:i}=e;if(ie(n),r&&!i){if(!n){g.value=A;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;m.value.some(e=>e[i]===r[i]||e[a]===r[a])||h.value.some(e=>e[i]===r[i]||e[a]===r[a])?g.value=A:g.value=[r]}}function De(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&X(),r&&!i&&a&&(h.value=A),K(),n?W([],[]):W(null,null)}function ke(e){!re(e,`action`)&&!re(e,`empty`)&&!re(e,`header`)&&e.preventDefault()}function Ae(e){ae(e)}function je(t){if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!w.value?.isComposing){if(C.value){let t=D.value?.getPendingTmNode();t?we(t):e.filterable||(X(),Me())}else if(Y(),e.tag&&de.value){let t=g.value[0];if(t){let n=t[e.valueField],{value:r}=u;e.multiple&&Array.isArray(r)&&r.includes(n)||Te(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;C.value&&D.value?.prev();break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;C.value?D.value?.next():Y();break;case`Escape`:C.value&&(H(t),X()),w.value?.focus()}}function Me(){w.value?.focus()}function Ne(){w.value?.focusInput()}function Pe(){C.value&&T.value?.syncPosition()}oe(),F(L(e,`options`),oe);let Fe={focus:()=>{w.value?.focus()},focusInput:()=>{w.value?.focusInput()},blur:()=>{w.value?.blur()},blurInput:()=>{w.value?.blurInput()}},Ie=s(()=>{let{self:{menuBoxShadow:e}}=o.value;return{"--n-menu-box-shadow":e}}),Le=i?ee(`select`,void 0,Ie,e):void 0;return{...Fe,mergedStatus:U,mergedClsPrefix:t,mergedBordered:n,namespace:r,treeMate:b,isMounted:te(),triggerRef:w,menuRef:D,pattern:p,uncontrolledShow:S,mergedShow:C,adjustedTo:J(e),uncontrolledValue:c,mergedValue:u,followerRef:T,localizedPlaceholder:k,selectedOption:R,selectedOptions:I,mergedSize:B,mergedDisabled:V,focused:d,activeWithoutMenuOpen:de,inlineThemeDisabled:i,onTriggerInputFocus:fe,onTriggerInputBlur:me,handleTriggerOrMenuResize:Pe,handleMenuFocus:_e,handleMenuBlur:ye,handleMenuTabOut:be,handleTriggerClick:Q,handleToggle:we,handleDeleteOption:Te,handlePatternInput:$,handleClear:De,handleTriggerBlur:he,handleTriggerFocus:ge,handleKeydown:je,handleMenuAfterLeave:le,handleMenuClickOutside:xe,handleMenuScroll:Ae,handleMenuKeydown:je,handleMenuMousedown:ke,mergedTheme:o,cssVars:i?void 0:Ie,themeClass:Le?.themeClass,onRender:Le?.onRender}},render(){return N(),t(`div`,{class:r(`${this.mergedClsPrefix}-select`)},[o(q,null,{_:1,default:n(()=>[(N(),m(de,null,{_:1,default:n(()=>(N(),m(dt,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{_:1,arrow:n(()=>[this.$slots.arrow?.()])},8,`inlineThemeDisabled.status.inputProps.clsPrefix.showArrow.maxTagCount.ellipsisTagPopoverProps.bordered.active.pattern.placeholder.selectedOption.selectedOptions.multiple.renderTag.renderLabel.filterable.clearable.disabled.size.theme.labelField.valueField.themeOverrides.loading.focused.onClick.onDeleteOption.onPatternInput.onClear.onBlur.onFocus.onKeydown.onPatternBlur.onPatternFocus.onResize.ignoreComposition`.split(`.`))))})),(N(),m(le,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===J.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{_:1,default:n(()=>(N(),m(i,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{_:1,default:n(()=>this.mergedShow||this.displayDirective===`show`?(this.onRender?.(),I((N(),m(Ze,c(this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{_:1,empty:n(()=>[this.$slots.empty?.()]),header:n(()=>[this.$slots.header?.()]),action:n(()=>[this.$slots.action?.()])},16,`onResize.inlineThemeDisabled.virtualScroll.class.clsPrefix.labelField.valueField.nodeProps.theme.themeOverrides.treeMate.multiple.size.renderOption.renderLabel.value.style.onToggle.onScroll.onFocus.onBlur.onKeydown.onTabOut.onMousedown.show.showCheckmark.resetMenuOnOptionsChange.scrollbarProps`.split(`.`))),this.displayDirective===`show`?[[A,this.mergedShow],[Se,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[Se,this.handleMenuClickOutside,void 0,{capture:!0}]])):null)},8,[`appear`,`onAfterLeave`])))},8,[`show`,`to`,`teleportDisabled`,`containerClass`,`width`,`placement`]))])})],2)}});export{ze as a,We as i,tt as n,Ze as r,mt as t};