import{$t as e,D as t,E as n,Et as r,Kt as i,Lt as a,Nt as o,P as s,Rt as c,Sn as l,T as u,Ut as d,Vt as f,Wt as p,Xt as m,Yt as h,Zt as g,a as _,bt as v,cn as y,d as b,dn as x,en as S,f as C,gn as w,gt as T,h as E,ht as D,in as O,j as k,k as A,kt as ee,l as j,ln as M,m as N,qt as P,rn as F,tn as I,u as te,vt as L,xt as R,yn as z,yt as B,zt as V}from"./vue-i18n-pPSPCP6m.js";import{c as ne,i as H,n as re,o as U,s as W}from"./event-BT3n3Ec6.js";import{c as G,d as K,f as q,g as J,h as ie,l as ae,m as oe,n as se,o as Y,p as ce,r as X,s as le,t as ue,u as de}from"./create-B2P7Jlzf.js";import{D as fe,F as pe,N as me,O as Z,P as Q,T as he,g as ge,j as _e,m as ve,v as ye,x as be,y as xe}from"./_plugin-vue_export-helper-cUzcR9Kp.js";import{a as Se,f as Ce,t as we,x as Te}from"./fade-in-scale-up.cssr-BU07VhSR.js";import{C as Ee,E as $,F as De,b as Oe}from"./index-BNdp-2__.js";function ke(e){return e&-e}var Ae=class{constructor(e,t){this.l=e,this.min=t;let n=Array(e+1);for(let t=0;t<e+1;++t)n[t]=0;this.ft=n}add(e,t){if(t===0)return;let{l:n,ft:r}=this;for(e+=1;e<=n;)r[e]+=t,e+=ke(e)}get(e){return this.sum(e+1)-this.sum(e)}sum(e){if(e===void 0&&(e=this.l),e<=0)return 0;let{ft:t,min:n,l:r}=this;if(e>r)throw Error("[FinweckTree.sum]: `i` is larger than length.");let i=e*n;for(;e>0;)i+=t[e],e-=ke(e);return i}getBound(e){let t=0,n=this.l;for(;n>t;){let r=Math.floor((t+n)/2),i=this.sum(r);if(i>e){n=r;continue}if(i<e){if(t===r)return this.sum(t+1)<=e?t+1:r;t=r}else return r}return t}},je;function Me(){return typeof document>`u`?!1:(je===void 0&&(je=`matchMedia`in window&&window.matchMedia(`(pointer:coarse)`).matches),je)}var Ne;function Pe(){return typeof document>`u`?1:(Ne===void 0&&(Ne=`chrome`in window?window.devicePixelRatio:1),Ne)}var Fe=`VVirtualListXScroll`;function Ie({columnsRef:e,renderColRef:t,renderItemWithColsRef:n}){let r=w(0),i=w(0),o=a(()=>{let t=e.value;if(t.length===0)return null;let n=new Ae(t.length,0);return t.forEach((e,t)=>{n.add(t,e.width)}),n}),s=b(()=>{let e=o.value;return e===null?0:Math.max(e.getBound(i.value)-1,0)}),c=e=>{let t=o.value;return t===null?0:t.sum(e)},l=b(()=>{let t=o.value;return t===null?0:Math.min(t.getBound(i.value+r.value)+1,e.value.length-1)});return O(Fe,{startIndexRef:s,endIndexRef:l,columnsRef:e,renderColRef:t,renderItemWithColsRef:n,getLeft:c}),{listWidthRef:r,scrollLeftRef:i}}var Le=p({name:`VirtualListRow`,props:{index:{type:Number,required:!0},item:{type:Object,required:!0}},setup(){let{startIndexRef:e,endIndexRef:t,columnsRef:n,getLeft:r,renderColRef:i,renderItemWithColsRef:a}=P(Fe);return{startIndex:e,endIndex:t,columns:n,renderCol:i,renderItemWithCols:a,getLeft:r}},render(){let{startIndex:e,endIndex:t,columns:n,renderCol:r,renderItemWithCols:i,getLeft:a,item:o}=this;if(i!=null)return i({itemIndex:this.index,startColIndex:e,endColIndex:t,allColumns:n,item:o,getLeft:a});if(r!=null){let i=[];for(let s=e;s<=t;++s){let e=n[s];i.push(r({column:e,left:a(s),item:o}))}return i}return null}}),Re=G(`.v-vl`,{maxHeight:`inherit`,height:`100%`,overflow:`auto`,minWidth:`1px`},[G(`&:not(.v-vl--show-scrollbar)`,{scrollbarWidth:`none`},[G(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,{width:0,height:0,display:`none`})])]),ze=p({name:`VirtualList`,inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},columns:{type:Array,default:()=>[]},renderCol:Function,renderItemWithCols:Function,items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:`div`},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:`key`},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(e){let t=k();Re.mount({id:`vueuc/virtual-list`,head:!0,anchorMetaName:ae,ssr:t}),I(()=>{let{defaultScrollIndex:t,defaultScrollKey:n}=e;t==null?n!=null&&v({key:n}):v({index:t})});let n=!1,r=!1;g(()=>{if(n=!1,!r){r=!0;return}v({top:m.value,left:s.value})}),S(()=>{n=!0,r||=!0});let i=b(()=>{if(e.renderCol==null&&e.renderItemWithCols==null||e.columns.length===0)return;let t=0;return e.columns.forEach(e=>{t+=e.width}),t}),o=a(()=>{let t=new Map,{keyField:n}=e;return e.items.forEach((e,r)=>{t.set(e[n],r)}),t}),{scrollLeftRef:s,listWidthRef:c}=Ie({columnsRef:z(e,`columns`),renderColRef:z(e,`renderCol`),renderItemWithColsRef:z(e,`renderItemWithCols`)}),l=w(null),u=w(void 0),d=new Map,f=a(()=>{let{items:t,itemSize:n,keyField:r}=e,i=new Ae(t.length,n);return t.forEach((e,t)=>{let n=e[r],a=d.get(n);a!==void 0&&i.add(t,a)}),i}),p=w(0),m=w(0),h=b(()=>Math.max(f.value.getBound(m.value-_e(e.paddingTop))-1,0)),_=a(()=>{let{value:t}=u;if(t===void 0)return[];let{items:n,itemSize:r}=e,i=h.value,a=Math.min(i+Math.ceil(t/r+1),n.length-1),o=[];for(let e=i;e<=a;++e)o.push(n[e]);return o}),v=(e,t)=>{if(typeof e==`number`){T(e,t,`auto`);return}let{left:n,top:r,index:i,key:a,position:s,behavior:c,debounce:l=!0}=e;if(n!==void 0||r!==void 0)T(n,r,c);else if(i!==void 0)C(i,c,l);else if(a!==void 0){let e=o.value.get(a);e!==void 0&&C(e,c,l)}else s===`bottom`?T(0,2**53-1,c):s===`top`&&T(0,0,c)},y,x=null;function C(t,n,r){let i=l.value;if(i==null)return;let{value:a}=f,o=a.sum(t)+_e(e.paddingTop);if(!r)i.scrollTo({left:0,top:o,behavior:n});else{y=t,x!==null&&window.clearTimeout(x),x=window.setTimeout(()=>{y=void 0,x=null},16);let{scrollTop:e,offsetHeight:r}=i;if(o>e){let s=a.get(t);o+s<=e+r||i.scrollTo({left:0,top:o+s-r,behavior:n})}else i.scrollTo({left:0,top:o,behavior:n})}}function T(e,t,n){l.value?.scrollTo({left:e,top:t,behavior:n})}function E(t,r){if(n||e.ignoreItemResize||N(r.target))return;let{value:i}=f,a=o.value.get(t),s=i.get(a),c=r.borderBoxSize?.[0]?.blockSize??r.contentRect.height;if(c===s)return;c-e.itemSize===0?d.delete(t):d.set(t,c-e.itemSize);let u=c-s;if(u===0)return;i.add(a,u);let m=l.value;if(m!=null){if(y===void 0){let e=i.sum(a);m.scrollTop>e&&m.scrollBy(0,u)}else(a<y||a===y&&c+i.sum(a)>m.scrollTop+m.offsetHeight)&&m.scrollBy(0,u);M()}p.value++}let D=!Me(),O=!1;function A(t){var n;(n=e.onScroll)==null||n.call(e,t),(!D||!O)&&M()}function ee(t){var n;if((n=e.onWheel)==null||n.call(e,t),D){let e=l.value;if(e!=null){if(t.deltaX===0&&(e.scrollTop===0&&t.deltaY<=0||e.scrollTop+e.offsetHeight>=e.scrollHeight&&t.deltaY>=0))return;t.preventDefault(),e.scrollTop+=t.deltaY/Pe(),e.scrollLeft+=t.deltaX/Pe(),M(),O=!0,Te(()=>{O=!1})}}}function j(t){if(n||N(t.target))return;if(e.renderCol==null&&e.renderItemWithCols==null){if(t.contentRect.height===u.value)return}else if(t.contentRect.height===u.value&&t.contentRect.width===c.value)return;u.value=t.contentRect.height,c.value=t.contentRect.width;let{onResize:r}=e;r!==void 0&&r(t)}function M(){let{value:e}=l;e!=null&&(m.value=e.scrollTop,s.value=e.scrollLeft)}function N(e){let t=e;for(;t!==null;){if(t.style.display===`none`)return!0;t=t.parentElement}return!1}return{listHeight:u,listStyle:{overflow:`auto`},keyToIndex:o,itemsStyle:a(()=>{let{itemResizable:t}=e,n=Q(f.value.sum());return p.value,[e.itemsStyle,{boxSizing:`content-box`,width:Q(i.value),height:t?``:n,minHeight:t?n:``,paddingTop:Q(e.paddingTop),paddingBottom:Q(e.paddingBottom)}]}),visibleItemsStyle:a(()=>(p.value,{transform:`translateY(${Q(f.value.sum(h.value))})`})),viewportItems:_,listElRef:l,itemsElRef:w(null),scrollTo:v,handleListResize:j,handleListScroll:A,handleListWheel:ee,handleItemResize:E}},render(){let{itemResizable:e,keyField:t,keyToIndex:n,visibleItemsTag:r}=this;return i(ye,{onResize:this.handleListResize},{default:()=>{var a;return i(`div`,h(this.$attrs,{class:[`v-vl`,this.showScrollbar&&`v-vl--show-scrollbar`],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:`listElRef`}),[this.items.length===0?(a=this.$slots).empty?.call(a):i(`div`,{ref:`itemsElRef`,class:`v-vl-items`,style:this.itemsStyle},[i(r,Object.assign({class:`v-vl-visible-items`,style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>{let{renderCol:r,renderItemWithCols:a}=this;return this.viewportItems.map(o=>{let s=o[t],c=n.get(s),l=r==null?void 0:i(Le,{index:c,item:o}),u=a==null?void 0:i(Le,{index:c,item:o}),d=this.$slots.default({item:o,renderedCols:l,renderedItemWithCols:u,index:c})[0];return e?i(ye,{key:s,onResize:e=>this.handleItemResize(s,e)},{default:()=>d}):(d.key=s,d)})}})])])}})}});function Be(e){switch(typeof e){case`string`:return e||void 0;case`number`:return String(e);default:return}}function Ve(t,n){n&&(I(()=>{let{value:e}=t;e&&xe.registerHandler(e,n)}),y(t,(e,t)=>{t&&xe.unregisterHandler(t)},{deep:!1}),e(()=>{let{value:e}=t;e&&xe.unregisterHandler(e)}))}var He=p({props:{onFocus:Function,onBlur:Function},setup(e){return()=>(()=>{let t=u(`d16ead82505dc285`);return F(),f(`div`,{style:`width: 0; height: 0`,tabindex:0,onFocus:t[0]||=t=>e.onFocus?.(t),onBlur:t[1]||=t=>e.onBlur?.(t)},null,32)})()}}),Ue=p({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=P(ie);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:r,tmNode:{rawNode:i}}=this,a=r?.(i),o=t?t(i,!1):$(i[this.labelField],i,!1),s=(F(),f(`div`,h(a,{class:[`${e}-base-select-group-header`,a?.class]}),[A(()=>o)],16));return i.render?i.render({node:s,option:i}):n?n({node:s,option:i,selected:!1}):s}});function We(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var Ge=p({name:`Checkmark`,render(){return(()=>{let e=u(`3c84eac8ae4e1f96`);return e[0]||=c(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},[c(`g`,{fill:`none`},[c(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})])],-1)})()}}),Ke=[`onClick`,`onMouseenter`,`onMousemove`];function qe(e,t){return F(),V(r,{name:`fade-in-scale-up-transition`},{default:()=>e?(F(),V(C,{key:1,clsPrefix:t,class:n(`${t}-base-select-option__check`)},{default:()=>i(Ge)},1032,[`clsPrefix`,`class`])):null},1024)}var Je=p({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:f}=P(ie),p=b(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function m(t){let{tmNode:n}=e;n.disabled||d(t,n)}function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function g(t){let{tmNode:n}=e,{value:r}=p;n.disabled||r||f(t,n)}return{multiple:r,isGrouped:b(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:p,isSelected:b(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:g,handleMouseEnter:h,handleClick:m}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:r,isPending:i,isGrouped:a,showCheckmark:o,nodeProps:s,renderOption:l,renderLabel:u,handleClick:d,handleMouseEnter:p,handleMouseMove:m}=this,g=qe(r,e),_=u?[u(t,r),o&&g]:[$(t[this.labelField],t,r),o&&g],v=s?.(t),y=(F(),f(`div`,h(v,{class:[`${e}-base-select-option`,t.class,v?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:r,[`${e}-base-select-option--grouped`]:a,[`${e}-base-select-option--pending`]:i,[`${e}-base-select-option--show-checkmark`]:o}],style:[v?.style||``,t.style||``],onClick:We([d,v?.onClick]),onMouseenter:We([p,v?.onMouseenter]),onMousemove:We([m,v?.onMousemove])}),[c(`div`,{class:n(`${e}-base-select-option__content`)},[A(()=>_)],2)],16,Ke));return t.render?t.render({node:y,option:t,selected:r}):l?l({node:y,option:t,selected:r}):y}}),Ye=T(`base-select-menu`,`
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
 `,[L(`content`,`
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
 `),L(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),L(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),L(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),L(`action`,`
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
 `,[B(`show-checkmark`,`
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
 `),B(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),B(`pending`,[D(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),B(`selected`,`
 color: var(--n-option-text-color-active);
 `,[D(`&::before`,`
 background-color: var(--n-option-color-active);
 `),B(`pending`,[D(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),B(`disabled`,`
 cursor: not-allowed;
 `,[v(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),B(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),L(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[we({enterScale:`0.5`})])])]),Xe=[`tabindex`,`onFocusin`,`onFocusout`,`onKeyup`,`onKeydown`,`onMousedown`,`onMouseenter`,`onMouseleave`],Ze=p({name:`InternalSelectMenu`,props:{...N.props,clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function},setup(t){let{mergedClsPrefixRef:n,mergedRtlRef:r,mergedComponentPropsRef:i}=s(t),o=j(`InternalSelectMenu`,r,n),c=N(`InternalSelectMenu`,`-internal-select-menu`,Ye,De,t,z(t,`clsPrefix`)),l=w(null),u=w(null),d=w(null),f=a(()=>t.treeMate.getFlattenedNodes()),p=a(()=>se(f.value)),h=w(null);function g(){let{treeMate:e}=t,n=null,{value:r}=t;r===null?n=e.getFirstAvailableNode():(n=t.multiple?e.getNode((r||[])[(r||[]).length-1]):e.getNode(r),(!n||n.disabled)&&(n=e.getFirstAvailableNode())),H(n||null)}function _(){let{value:e}=h;e&&!t.treeMate.getNode(e.key)&&(h.value=null)}let v;y(()=>t.show,e=>{e?v=y(()=>t.treeMate,()=>{t.resetMenuOnOptionsChange?(t.autoPending?g():_(),m(re)):_()},{immediate:!0}):v?.()},{immediate:!0}),e(()=>{v?.()});let b=a(()=>_e(c.value.self[R(`optionHeight`,t.size)])),x=a(()=>me(c.value.self[R(`padding`,t.size)])),S=a(()=>t.multiple&&Array.isArray(t.value)?new Set(t.value):new Set),C=a(()=>{let e=f.value;return e&&e.length===0}),T=a(()=>i?.value?.Select?.renderEmpty);function D(e){let{onToggle:n}=t;n&&n(e)}function k(e){let{onScroll:n}=t;n&&n(e)}function A(e){d.value?.sync(),k(e)}function ee(){d.value?.sync()}function M(){let{value:e}=h;return e||null}function P(e,t){t.disabled||H(t,!1)}function F(e,t){t.disabled||D(t)}function te(e){J(e,`action`)||t.onKeyup?.(e)}function L(e){J(e,`action`)||t.onKeydown?.(e)}function B(e){t.onMousedown?.(e),!t.focusable&&e.preventDefault()}function V(){let{value:e}=h;e&&H(e.getNext({loop:!0}),!0)}function ne(){let{value:e}=h;e&&H(e.getPrev({loop:!0}),!0)}function H(e,t=!1){h.value=e,t&&re()}function re(){let e=h.value;if(!e)return;let n=p.value(e.key);n!==null&&(t.virtualScroll?u.value?.scrollTo({index:n}):d.value?.scrollTo({index:n,elSize:b.value}))}function U(e){l.value?.contains(e.target)&&t.onFocus?.(e)}function W(e){l.value?.contains(e.relatedTarget)||t.onBlur?.(e)}O(ie,{handleOptionMouseEnter:P,handleOptionClick:F,valueSetRef:S,pendingTmNodeRef:h,nodePropsRef:z(t,`nodeProps`),showCheckmarkRef:z(t,`showCheckmark`),multipleRef:z(t,`multiple`),valueRef:z(t,`value`),renderLabelRef:z(t,`renderLabel`),renderOptionRef:z(t,`renderOption`),labelFieldRef:z(t,`labelField`),valueFieldRef:z(t,`valueField`)}),O(oe,l),I(()=>{let{value:e}=d;e&&e.sync()});let G=a(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:a,groupHeaderTextColor:o,actionDividerColor:s,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[R(`optionFontSize`,e)]:x,[R(`optionHeight`,e)]:S,[R(`optionPadding`,e)]:C}}=c.value;return{"--n-height":r,"--n-action-divider-color":s,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":a,"--n-option-font-size":x,"--n-group-header-text-color":o,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":me(C,`left`),"--n-option-padding-right":me(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:K}=t,q=K?E(`internal-select-menu`,a(()=>t.size[0]),G,t):void 0,ae={selfRef:l,next:V,prev:ne,getPendingTmNode:M};return Ve(l,t.onResize),{mergedTheme:c,mergedClsPrefix:n,rtlEnabled:o,virtualListRef:u,scrollbarRef:d,itemSize:b,padding:x,flattenedNodes:f,empty:C,mergedRenderEmpty:T,virtualListContainer(){let{value:e}=u;return e?.listElRef},virtualListContent(){let{value:e}=u;return e?.itemsElRef},doScroll:k,handleFocusin:U,handleFocusout:W,handleKeyUp:te,handleKeyDown:L,handleMouseDown:B,handleVirtualListResize:ee,handleVirtualListScroll:A,cssVars:K?void 0:G,themeClass:q?.themeClass,onRender:q?.onRender,...ae}},render(){let{$slots:e,virtualScroll:t,clsPrefix:r,mergedTheme:i,themeClass:a,onRender:s}=this;return s?.(),F(),f(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:n([`${r}-base-select-menu`,`${r}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${r}-base-select-menu--rtl`,a,this.multiple&&`${r}-base-select-menu--multiple`]),style:l(this.cssVars),onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},[A(()=>fe(e.header,e=>e&&(F(),f(`div`,{class:n(`${r}-base-select-menu__header`),"data-header":!0,key:`header`},[A(()=>e)],2)))),this.loading?(F(),f(`div`,{key:0,class:n(`${r}-base-select-menu__loading`)},[(F(),V(_,{clsPrefix:r,strokeWidth:20},null,8,[`clsPrefix`]))],2)):(F(),f(o,{key:1},[this.empty?(F(),f(`div`,{key:1,class:n(`${r}-base-select-menu__empty`),"data-empty":!0},[A(()=>he(e.empty,()=>[this.mergedRenderEmpty?.()||(F(),V(W,{theme:i.peers.Empty,themeOverrides:i.peerOverrides.Empty,size:this.size},null,8,[`theme`,`themeOverrides`,`size`]))]))],2)):(F(),V(ge,h({key:0,ref:`scrollbarRef`,theme:i.peers.Scrollbar,themeOverrides:i.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?(F(),V(ze,{key:1,ref:`virtualListRef`,class:n(`${r}-virtual-list`),items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?(F(),V(Ue,{key:e.key,clsPrefix:r,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):e.ignored?null:(F(),V(Je,{clsPrefix:r,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))},1032,[`class`,`items`,`itemSize`,`paddingTop`,`paddingBottom`,`onResize`,`onScroll`])):(F(),f(`div`,{key:4,class:n(`${r}-base-select-menu-option-wrapper`),style:l({paddingTop:this.padding.top,paddingBottom:this.padding.bottom})},[A(()=>this.flattenedNodes.map(e=>e.isGroup?(F(),V(Ue,{key:e.key,clsPrefix:r,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):(F(),V(Je,{clsPrefix:r,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))))],6))},1040,[`theme`,`themeOverrides`,`scrollable`,`container`,`content`,`onScroll`]))],64)),A(()=>fe(e.action,e=>e&&[(F(),f(`div`,{class:n(`${r}-base-select-menu__action`),"data-action":!0,key:`action`},[A(()=>e)],2)),(F(),V(He,{onFocus:this.onTabOut,key:`focus-detector`},null,8,[`onFocus`]))]))],46,Xe)}});function Qe(e){return e.type===`group`}function $e(e){return e.type===`ignored`}function et(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function tt(e,t){return{getIsGroup:Qe,getIgnored:$e,getKey(t){return Qe(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function nt(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(Qe(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if($e(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function rt(e,t,n){let r=new Map;return e.forEach(e=>{Qe(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var it=D([T(`base-selection`,`
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
 `),T(`base-selection-tags`,`min-height: var(--n-height);`),L(`border, state-border`,`
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
 `),L(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),T(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[L(`arrow`,`
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
 `,[L(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),T(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[L(`inner`,`
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
 `,[L(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),L(`render-label`,`
 color: var(--n-text-color);
 `)]),v(`disabled`,[D(`&:hover`,[L(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),B(`focus`,[L(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),B(`active`,[L(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),T(`base-selection-label`,`background-color: var(--n-color-active);`),T(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),B(`disabled`,`cursor: not-allowed;`,[L(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),T(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[T(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),L(`render-label`,`
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
 `,[L(`input`,`
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
 `),L(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>B(`${e}-status`,[L(`state-border`,`border: var(--n-border-${e});`),v(`disabled`,[D(`&:hover`,[L(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),B(`active`,[L(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),T(`base-selection-label`,`background-color: var(--n-color-active-${e});`),T(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),B(`focus`,[L(`state-border`,`
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
 `,[L(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),at=[`disabled`,`value`,`autofocus`,`onBlur`,`onFocus`,`onKeydown`,`onInput`,`onCompositionstart`,`onCompositionend`],ot=[`tabindex`],st=[`title`],ct=[`value`,`readonly`,`disabled`,`autofocus`,`onFocus`,`onBlur`,`onInput`,`onCompositionstart`,`onCompositionend`],lt=[`tabindex`],ut=[`onClick`,`onMouseenter`,`onMouseleave`,`onKeydown`,`onFocusin`,`onFocusout`,`onMousedown`],dt=p({name:`InternalSelection`,props:{...N.props,clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function},setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n}=s(e),r=j(`InternalSelection`,n,t),i=w(null),o=w(null),c=w(null),l=w(null),u=w(null),d=w(null),f=w(null),p=w(null),h=w(null),g=w(null),_=w(!1),v=w(!1),b=w(!1),x=N(`InternalSelection`,`-internal-selection`,it,Ee,e,z(e,`clsPrefix`)),S=a(()=>e.clearable&&!e.disabled&&(b.value||e.active)),C=a(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):$(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),T=a(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),D=a(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function O(){let{value:t}=i;if(t){let{value:n}=o;n&&(n.style.width=`${t.offsetWidth}px`,e.maxTagCount!==`responsive`&&h.value?.sync({showAllItemsBeforeCalculate:!1}))}}function k(){let{value:e}=g;e&&(e.style.display=`none`)}function A(){let{value:e}=g;e&&(e.style.display=`inline-block`)}y(z(e,`active`),e=>{e||k()}),y(z(e,`pattern`),()=>{e.multiple&&m(O)});function ee(t){let{onFocus:n}=e;n&&n(t)}function P(t){let{onBlur:n}=e;n&&n(t)}function F(t){let{onDeleteOption:n}=e;n&&n(t)}function te(t){let{onClear:n}=e;n&&n(t)}function L(t){let{onPatternInput:n}=e;n&&n(t)}function B(e){(!e.relatedTarget||!c.value?.contains(e.relatedTarget))&&ee(e)}function V(e){c.value?.contains(e.relatedTarget)||P(e)}function ne(e){te(e)}function H(){b.value=!0}function re(){b.value=!1}function U(t){e.active&&e.filterable&&t.target!==o.value&&t.preventDefault()}function W(e){F(e)}let G=w(!1);function K(t){if(t.key===`Backspace`&&!G.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&W(t[t.length-1])}}let q=null;function J(t){let{value:n}=i;n&&(n.textContent=t.target.value,O()),e.ignoreComposition&&G.value?q=t:L(t)}function ie(){G.value=!0}function ae(){G.value=!1,e.ignoreComposition&&L(q),q=null}function oe(t){v.value=!0,e.onPatternFocus?.(t)}function se(t){v.value=!1,e.onPatternBlur?.(t)}function Y(){if(e.filterable)v.value=!1,d.value?.blur(),o.value?.blur();else if(e.multiple){let{value:e}=l;e?.blur()}else{let{value:e}=u;e?.blur()}}function ce(){e.filterable?(v.value=!1,d.value?.focus()):e.multiple?l.value?.focus():u.value?.focus()}function X(){let{value:e}=o;e&&(A(),e.focus())}function le(){let{value:e}=o;e&&e.blur()}function ue(e){let{value:t}=f;t&&t.setTextContent(`+${e}`)}function de(){let{value:e}=p;return e}function fe(){return o.value}let pe=null;function Z(){pe!==null&&window.clearTimeout(pe)}function Q(){e.active||(Z(),pe=window.setTimeout(()=>{D.value&&(_.value=!0)},100))}function he(){Z()}function ge(e){e||(Z(),_.value=!1)}y(D,e=>{e||(_.value=!1)}),I(()=>{M(()=>{let t=d.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=v.value?-1:0)})}),Ve(c,e.onResize);let{inlineThemeDisabled:_e}=e,ve=a(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:b,borderActive:S,arrowColor:C,arrowColorDisabled:w,loadingColor:T,colorActiveWarning:E,boxShadowFocusWarning:D,boxShadowActiveWarning:O,boxShadowHoverWarning:k,borderWarning:A,borderFocusWarning:ee,borderHoverWarning:j,borderActiveWarning:M,colorActiveError:N,boxShadowFocusError:P,boxShadowActiveError:F,boxShadowHoverError:I,borderError:te,borderFocusError:L,borderHoverError:z,borderActiveError:B,clearColor:V,clearColorHover:ne,clearColorPressed:H,clearSize:re,arrowSize:U,[R(`height`,t)]:W,[R(`fontSize`,t)]:G}}=x.value,K=me(c),q=me(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":y,"--n-border-hover":b,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":G,"--n-height":W,"--n-padding-single-top":K.top,"--n-padding-multiple-top":q.top,"--n-padding-single-right":K.right,"--n-padding-multiple-right":q.right,"--n-padding-single-left":K.left,"--n-padding-multiple-left":q.left,"--n-padding-single-bottom":K.bottom,"--n-padding-multiple-bottom":q.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":C,"--n-arrow-color-disabled":w,"--n-loading-color":T,"--n-color-active-warning":E,"--n-box-shadow-focus-warning":D,"--n-box-shadow-active-warning":O,"--n-box-shadow-hover-warning":k,"--n-border-warning":A,"--n-border-focus-warning":ee,"--n-border-hover-warning":j,"--n-border-active-warning":M,"--n-color-active-error":N,"--n-box-shadow-focus-error":P,"--n-box-shadow-active-error":F,"--n-box-shadow-hover-error":I,"--n-border-error":te,"--n-border-focus-error":L,"--n-border-hover-error":z,"--n-border-active-error":B,"--n-clear-size":re,"--n-clear-color":V,"--n-clear-color-hover":ne,"--n-clear-color-pressed":H,"--n-arrow-size":U,"--n-font-weight":r}}),ye=_e?E(`internal-selection`,a(()=>e.size[0]),ve,e):void 0;return{mergedTheme:x,mergedClearable:S,mergedClsPrefix:t,rtlEnabled:r,patternInputFocused:v,filterablePlaceholder:C,label:T,selected:D,showTagsPanel:_,isComposing:G,counterRef:f,counterWrapperRef:p,patternInputMirrorRef:i,patternInputRef:o,selfRef:c,multipleElRef:l,singleElRef:u,patternInputWrapperRef:d,overflowRef:h,inputTagElRef:g,handleMouseDown:U,handleFocusin:B,handleClear:ne,handleMouseEnter:H,handleMouseLeave:re,handleDeleteOption:W,handlePatternKeyDown:K,handlePatternInputInput:J,handlePatternInputBlur:se,handlePatternInputFocus:oe,handleMouseEnterCounter:Q,handleMouseLeaveCounter:he,handleFocusout:V,handleCompositionEnd:ae,handleCompositionStart:ie,onPopoverUpdateShow:ge,focus:ce,focusInput:X,blur:Y,blurInput:le,updateCounter:ue,getCounter:de,getTail:fe,renderLabel:e.renderLabel,cssVars:_e?void 0:ve,themeClass:ye?.themeClass,onRender:ye?.onRender}},render(){let{status:e,multiple:t,size:r,disabled:i,filterable:a,maxTagCount:s,bordered:u,clsPrefix:d,ellipsisTagPopoverProps:p,onRender:m,renderTag:g,renderLabel:_}=this;m?.();let v=s===`responsive`,y=typeof s==`number`,b=v||y,x=(F(),V(be,null,{default:()=>(F(),V(H,{clsPrefix:d,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>this.$slots.arrow?.()},1032,[`clsPrefix`,`loading`,`showArrow`,`showClear`,`onClear`]))},1024)),S;if(t){let{labelField:e}=this,t=t=>(F(),f(`div`,{class:n(`${d}-base-selection-tag-wrapper`),key:t.value},[g?(F(),f(o,{key:0},[A(()=>g({option:t,handleClose:()=>{this.handleDeleteOption(t)}}))],64)):(F(),V(U,{key:1,size:r,closable:!t.disabled,disabled:i,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>_?_(t,!0):$(t[e],t,!0)},1032,[`size`,`closable`,`disabled`,`onClose`]))],2)),l=()=>(y?this.selectedOptions.slice(0,s):this.selectedOptions).map(t),u=a?(F(),f(`div`,{class:n(`${d}-base-selection-input-tag`),ref:`inputTagElRef`,key:`__input-tag__`},[c(`input`,h(this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:i,value:this.pattern,autofocus:this.autofocus,class:`${d}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,at),c(`span`,{ref:`patternInputMirrorRef`,class:n(`${d}-base-selection-input-tag__mirror`)},[A(()=>this.pattern)],2)],2)):null,m=v?()=>(F(),f(`div`,{class:n(`${d}-base-selection-tag-wrapper`),ref:`counterWrapperRef`},[(F(),V(U,{size:r,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:i},null,8,[`size`,`onMouseenter`,`onMouseleave`,`disabled`]))],2)):void 0,C;if(y){let e=this.selectedOptions.length-s;e>0&&(C=(t=>(F(),f(`div`,{class:n(`${d}-base-selection-tag-wrapper`),key:`__counter__`},[(F(),V(U,{size:r,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:i},{default:()=>`+${e}`},1032,[`size`,`onMouseenter`,`disabled`]))],2)))(C))}let w=v?a?(F(),V(Y,{key:3,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:l,counter:m,tail:()=>u},1032,[`updateCounter`,`getCounter`,`getTail`])):(F(),V(Y,{key:4,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:l,counter:m},1032,[`updateCounter`,`getCounter`])):y&&C?l().concat(C):l(),T=b?()=>(F(),f(`div`,{class:n(`${d}-base-selection-popover`)},[v?(F(),f(o,{key:0},[A(()=>l())],64)):(F(),f(o,{key:1},[A(()=>this.selectedOptions.map(t))],64))],2)):void 0,E=b?{show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover,...p}:null,D=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?(F(),f(`div`,{key:5,class:n(`${d}-base-selection-placeholder ${d}-base-selection-overlay`)},[c(`div`,{class:n(`${d}-base-selection-placeholder__inner`)},[A(()=>this.placeholder)],2)],2)):null,O=a?(F(),f(`div`,{key:6,ref:`patternInputWrapperRef`,class:n(`${d}-base-selection-tags`)},[A(()=>w),v?A(()=>null):(F(),f(o,{key:1},[A(()=>u)],64)),A(()=>x)],2)):(F(),f(`div`,{key:7,ref:`multipleElRef`,class:n(`${d}-base-selection-tags`),tabindex:i?void 0:0},[A(()=>w),A(()=>x)],10,ot));S=(e=>(F(),f(o,{key:8},[b?(F(),V(X,h({key:0},E,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>O,default:T},1040)):(F(),f(o,{key:1},[A(()=>O)],64)),A(()=>D)],64)))(S)}else if(a){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,r=!this.active&&this.selected;S=(e=>(F(),f(`div`,{key:9,ref:`patternInputWrapperRef`,class:n(`${d}-base-selection-label`),title:this.patternInputFocused?void 0:Be(this.label)},[c(`input`,h(this.inputProps,{ref:`patternInputRef`,class:`${d}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:i,disabled:i,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,ct),r?(F(),f(`div`,{class:n(`${d}-base-selection-label__render-label ${d}-base-selection-overlay`),key:`input`},[c(`div`,{class:n(`${d}-base-selection-overlay__wrapper`)},[g?(F(),f(o,{key:0},[A(()=>g({option:this.selectedOption,handleClose:()=>{}}))],64)):(F(),f(o,{key:1},[_?(F(),f(o,{key:0},[A(()=>_(this.selectedOption,!0))],64)):(F(),f(o,{key:1},[A(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],2)):A(()=>null),t?(F(),f(`div`,{class:n(`${d}-base-selection-placeholder ${d}-base-selection-overlay`),key:`placeholder`},[c(`div`,{class:n(`${d}-base-selection-overlay__wrapper`)},[A(()=>this.filterablePlaceholder)],2)],2)):A(()=>null),A(()=>x)],10,st)))(S)}else S=(e=>(F(),f(`div`,{key:10,ref:`singleElRef`,class:n(`${d}-base-selection-label`),tabindex:this.disabled?void 0:0},[this.label===void 0?(F(),f(`div`,{class:n(`${d}-base-selection-placeholder ${d}-base-selection-overlay`),key:`placeholder`},[c(`div`,{class:n(`${d}-base-selection-placeholder__inner`)},[A(()=>this.placeholder)],2)],2)):(F(),f(`div`,{class:n(`${d}-base-selection-input`),title:Be(this.label),key:`input`},[c(`div`,{class:n(`${d}-base-selection-input__content`)},[g?(F(),f(o,{key:0},[A(()=>g({option:this.selectedOption,handleClose:()=>{}}))],64)):(F(),f(o,{key:1},[_?(F(),f(o,{key:0},[A(()=>_(this.selectedOption,!0))],64)):(F(),f(o,{key:1},[A(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],10,[`title`])),A(()=>x)],10,lt)))(S);return F(),f(`div`,{ref:`selfRef`,class:n([`${d}-base-selection`,this.rtlEnabled&&`${d}-base-selection--rtl`,this.themeClass,e&&`${d}-base-selection--${e}-status`,{[`${d}-base-selection--active`]:this.active,[`${d}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${d}-base-selection--disabled`]:this.disabled,[`${d}-base-selection--multiple`]:this.multiple,[`${d}-base-selection--focus`]:this.focused}]),style:l(this.cssVars),onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},[A(()=>S),u?(F(),f(`div`,{key:0,class:n(`${d}-base-selection__border`)},null,2)):A(()=>null),u?(F(),f(`div`,{key:2,class:n(`${d}-base-selection__state-border`)},null,2)):A(()=>null)],46,ut)}}),ft=D([T(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),T(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[we({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),pt={...N.props,to:q.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array},mt=p({name:`Select`,props:pt,slots:Object,setup(e){let{mergedClsPrefixRef:t,mergedBorderedRef:n,namespaceRef:r,inlineThemeDisabled:i,mergedComponentPropsRef:o}=s(e),c=N(`Select`,`-select`,ft,Oe,e,t),l=w(e.defaultValue),u=z(e,`value`),d=Ce(u,l),f=w(!1),p=w(``),m=ce(e,[`items`,`options`]),h=w([]),g=w([]),_=a(()=>g.value.concat(h.value).concat(m.value)),v=a(()=>{let{filter:t}=e;if(t)return t;let{labelField:n,valueField:r}=e;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return et(e,i);let a=t[r];return typeof a==`string`?et(e,a):typeof a==`number`&&et(e,String(a))}}),b=a(()=>{if(e.remote)return m.value;{let{value:t}=_,{value:n}=p;return!n.length||!e.filterable?t:nt(t,v.value,n,e.childrenField)}}),x=a(()=>{let{valueField:t,childrenField:n}=e,r=tt(t,n);return ue(b.value,r)}),S=a(()=>rt(_.value,e.valueField,e.childrenField)),C=w(!1),T=Ce(z(e,`show`),C),D=w(null),O=w(null),k=w(null),{localeRef:A}=ne(`Select`),ee=a(()=>e.placeholder??A.value.placeholder),j=[],M=w(new Map),P=a(()=>{let{fallbackOption:t}=e;if(t===void 0){let{labelField:t,valueField:n}=e;return e=>({[t]:String(e),[n]:e})}return t===!1?!1:e=>Object.assign(t(e),{value:e})});function F(t){let n=e.remote,{value:r}=M,{value:i}=S,{value:a}=P,o=[];return t.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let I=a(()=>{if(e.multiple){let{value:e}=d;return Array.isArray(e)?F(e):[]}return null}),L=a(()=>{let{value:t}=d;return!e.multiple&&!Array.isArray(t)?t===null?null:F([t])[0]||null:null}),R=ve(e,{mergedSize:t=>{let{size:n}=e;if(n)return n;let{mergedSize:r}=t||{};return r?.value?r.value:o?.value?.Select?.size||`medium`}}),{mergedSizeRef:B,mergedDisabledRef:V,mergedStatusRef:H}=R;function U(t,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=e,{nTriggerFormChange:o,nTriggerFormInput:s}=R;r&&Z(r,t,n),a&&Z(a,t,n),i&&Z(i,t,n),l.value=t,o(),s()}function W(t){let{onBlur:n}=e,{nTriggerFormBlur:r}=R;n&&Z(n,t),r()}function G(){let{onClear:t}=e;t&&Z(t)}function K(t){let{onFocus:n,showOnFocus:r}=e,{nTriggerFormFocus:i}=R;n&&Z(n,t),i(),r&&Y()}function ie(t){let{onSearch:n}=e;n&&Z(n,t)}function ae(t){let{onScroll:n}=e;n&&Z(n,t)}function oe(){let{remote:t,multiple:n}=e;if(t){let{value:t}=M;if(n){let{valueField:n}=e;I.value?.forEach(e=>{t.set(e[n],e)})}else{let n=L.value;n&&t.set(n[e.valueField],n)}}}function se(t){let{onUpdateShow:n,"onUpdate:show":r}=e;n&&Z(n,t),r&&Z(r,t),C.value=t}function Y(){V.value||(se(!0),C.value=!0,e.filterable&&Ne())}function X(){se(!1)}function le(){p.value=``,g.value=j}let de=w(!1);function fe(){e.filterable&&(de.value=!0)}function me(){e.filterable&&(de.value=!1,T.value||le())}function Q(){V.value||(T.value?e.filterable?Ne():X():Y())}function he(e){k.value?.selfRef?.contains(e.relatedTarget)||(f.value=!1,W(e),X())}function ge(e){K(e),f.value=!0}function _e(){f.value=!0}function ye(e){D.value?.$el.contains(e.relatedTarget)||(f.value=!1,W(e),X())}function be(){D.value?.focus(),X()}function xe(e){T.value&&(D.value?.$el.contains(pe(e))||X())}function Se(t){if(!Array.isArray(t))return[];if(P.value)return Array.from(t);{let{remote:n}=e,{value:r}=S;if(n){let{value:e}=M;return t.filter(t=>r.has(t)||e.has(t))}return t.filter(e=>r.has(e))}}function we(e){Te(e.rawNode)}function Te(t){if(V.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=e;if(n&&!r){let{value:e}=g,t=e[0]||null;if(t){let e=h.value;e.length?e.push(t):h.value=[t],g.value=j}}if(r&&M.value.set(t[a],t),e.multiple){let e=Se(d.value),o=e.findIndex(e=>e===t[a]);if(~o){if(e.splice(o,1),n&&!r){let e=Ee(t[a]);~e&&(h.value.splice(e,1),i&&(p.value=``))}}else e.push(t[a]),i&&(p.value=``);U(e,F(e))}else{if(n&&!r){let e=Ee(t[a]);~e?h.value=[h.value[e]]:h.value=j}Me(),X(),U(t[a],t)}}function Ee(t){return h.value.findIndex(n=>n[e.valueField]===t)}function $(t){T.value||Y();let{value:n}=t.target;p.value=n;let{tag:r,remote:i}=e;if(ie(n),r&&!i){if(!n){g.value=j;return}let{onCreate:t}=e,r=t?t(n):{[e.labelField]:n,[e.valueField]:n},{valueField:i,labelField:a}=e;m.value.some(e=>e[i]===r[i]||e[a]===r[a])||h.value.some(e=>e[i]===r[i]||e[a]===r[a])?g.value=j:g.value=[r]}}function De(t){t.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=e;!n&&e.filterable&&X(),r&&!i&&a&&(h.value=j),G(),n?U([],[]):U(null,null)}function ke(e){!J(e,`action`)&&!J(e,`empty`)&&!J(e,`header`)&&e.preventDefault()}function Ae(e){ae(e)}function je(t){if(!e.keyboard){t.preventDefault();return}switch(t.key){case` `:if(e.filterable)break;t.preventDefault();case`Enter`:if(!D.value?.isComposing){if(T.value){let t=k.value?.getPendingTmNode();t?we(t):e.filterable||(X(),Me())}else if(Y(),e.tag&&de.value){let t=g.value[0];if(t){let n=t[e.valueField],{value:r}=d;e.multiple&&Array.isArray(r)&&r.includes(n)||Te(t)}}}t.preventDefault();break;case`ArrowUp`:if(t.preventDefault(),e.loading)return;T.value&&k.value?.prev();break;case`ArrowDown`:if(t.preventDefault(),e.loading)return;T.value?k.value?.next():Y();break;case`Escape`:T.value&&(re(t),X()),D.value?.focus()}}function Me(){D.value?.focus()}function Ne(){D.value?.focusInput()}function Pe(){T.value&&O.value?.syncPosition()}oe(),y(z(e,`options`),oe);let Fe={focus:()=>{D.value?.focus()},focusInput:()=>{D.value?.focusInput()},blur:()=>{D.value?.blur()},blurInput:()=>{D.value?.blurInput()}},Ie=a(()=>{let{self:{menuBoxShadow:e}}=c.value;return{"--n-menu-box-shadow":e}}),Le=i?E(`select`,void 0,Ie,e):void 0;return{...Fe,mergedStatus:H,mergedClsPrefix:t,mergedBordered:n,namespace:r,treeMate:x,isMounted:te(),triggerRef:D,menuRef:k,pattern:p,uncontrolledShow:C,mergedShow:T,adjustedTo:q(e),uncontrolledValue:l,mergedValue:d,followerRef:O,localizedPlaceholder:ee,selectedOption:L,selectedOptions:I,mergedSize:B,mergedDisabled:V,focused:f,activeWithoutMenuOpen:de,inlineThemeDisabled:i,onTriggerInputFocus:fe,onTriggerInputBlur:me,handleTriggerOrMenuResize:Pe,handleMenuFocus:_e,handleMenuBlur:ye,handleMenuTabOut:be,handleTriggerClick:Q,handleToggle:we,handleDeleteOption:Te,handlePatternInput:$,handleClear:De,handleTriggerBlur:he,handleTriggerFocus:ge,handleKeydown:je,handleMenuAfterLeave:le,handleMenuClickOutside:xe,handleMenuScroll:Ae,handleMenuKeydown:je,handleMenuMousedown:ke,mergedTheme:c,cssVars:i?void 0:Ie,themeClass:Le?.themeClass,onRender:Le?.onRender}},render(){return F(),f(`div`,{class:n(`${this.mergedClsPrefix}-select`)},[d(K,null,{_:1,default:t(()=>[(F(),V(de,null,{_:1,default:t(()=>(F(),V(dt,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{_:1,arrow:t(()=>[this.$slots.arrow?.()])},8,`inlineThemeDisabled.status.inputProps.clsPrefix.showArrow.maxTagCount.ellipsisTagPopoverProps.bordered.active.pattern.placeholder.selectedOption.selectedOptions.multiple.renderTag.renderLabel.filterable.clearable.disabled.size.theme.labelField.valueField.themeOverrides.loading.focused.onClick.onDeleteOption.onPatternInput.onClear.onBlur.onFocus.onKeydown.onPatternBlur.onPatternFocus.onResize.ignoreComposition`.split(`.`))))})),(F(),V(le,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===q.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{_:1,default:t(()=>(F(),V(r,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{_:1,default:t(()=>this.mergedShow||this.displayDirective===`show`?(this.onRender?.(),x((F(),V(Ze,h(this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{_:1,empty:t(()=>[this.$slots.empty?.()]),header:t(()=>[this.$slots.header?.()]),action:t(()=>[this.$slots.action?.()])},16,`onResize.inlineThemeDisabled.virtualScroll.class.clsPrefix.labelField.valueField.nodeProps.theme.themeOverrides.treeMate.multiple.size.renderOption.renderLabel.value.style.onToggle.onScroll.onFocus.onBlur.onKeydown.onTabOut.onMousedown.show.showCheckmark.resetMenuOnOptionsChange.scrollbarProps`.split(`.`))),this.displayDirective===`show`?[[ee,this.mergedShow],[Se,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[Se,this.handleMenuClickOutside,void 0,{capture:!0}]])):null)},8,[`appear`,`onAfterLeave`])))},8,[`show`,`to`,`teleportDisabled`,`containerClass`,`width`,`placement`]))])})],2)}});export{ze as a,We as i,tt as n,Ze as r,mt as t};