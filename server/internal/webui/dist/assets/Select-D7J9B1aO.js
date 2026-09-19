import{$t as e,At as t,Bt as n,Ft as r,Gt as i,Ht as a,Jt as o,Kt as s,Lt as c,Mt as l,Ot as u,Pt as d,Rt as f,Vt as p,Wt as m,Yt as h,_ as g,at as _,bt as v,ct as y,dt as b,en as x,g as S,gt as C,h as w,i as T,kt as E,lt as D,nn as O,o as k,ot as A,pn as j,r as M,s as ee,sn as N,un as P,ut as F,w as te,wt as I,x as ne,y as L}from"./vue-i18n-CgKceEWx.js";import{n as re,o as ie,r as R,s as ae}from"./event-BLioqdn-.js";import{c as z,d as B,f as oe,g as V,h as H,l as U,m as W,n as se,o as G,p as ce,r as K,s as le,t as ue,u as q}from"./create-Crm1UAxV.js";import{_ as J,a as Y,b as X,d as de,i as fe,m as Z,n as pe,p as me,s as he,x as ge,y as _e}from"./_plugin-vue_export-helper-BcoWp_Ef.js";import{a as ve,f as ye,t as be,x as xe}from"./fade-in-scale-up.cssr-7akLf3Ad.js";import{n as Se,r as Q}from"./FadeInExpandTransition-D8IUT6y9.js";import{t as Ce}from"./Tag-DcxjiX9Q.js";import{m as we}from"./http-BdILZoCd.js";import{C as Te,D as $,V as Ee,b as De,k as Oe}from"./index-kjlF3dpY.js";function ke(e){return e&-e}var Ae=class{constructor(e,t){this.l=e,this.min=t;let n=Array(e+1);for(let t=0;t<e+1;++t)n[t]=0;this.ft=n}add(e,t){if(t===0)return;let{l:n,ft:r}=this;for(e+=1;e<=n;)r[e]+=t,e+=ke(e)}get(e){return this.sum(e+1)-this.sum(e)}sum(e){if(e===void 0&&(e=this.l),e<=0)return 0;let{ft:t,min:n,l:r}=this;if(e>r)throw Error("[FinweckTree.sum]: `i` is larger than length.");let i=e*n;for(;e>0;)i+=t[e],e-=ke(e);return i}getBound(e){let t=0,n=this.l;for(;n>t;){let r=Math.floor((t+n)/2),i=this.sum(r);if(i>e){n=r;continue}if(i<e){if(t===r)return this.sum(t+1)<=e?t+1:r;t=r}else return r}return t}},je;function Me(){return typeof document>`u`?!1:(je===void 0&&(je=`matchMedia`in window&&window.matchMedia(`(pointer:coarse)`).matches),je)}var Ne;function Pe(){return typeof document>`u`?1:(Ne===void 0&&(Ne=`chrome`in window?window.devicePixelRatio:1),Ne)}var Fe=`VVirtualListXScroll`;function Ie({columnsRef:e,renderColRef:t,renderItemWithColsRef:n}){let r=N(0),i=N(0),a=u(()=>{let t=e.value;if(t.length===0)return null;let n=new Ae(t.length,0);return t.forEach((e,t)=>{n.add(t,e.width)}),n}),o=Q(()=>{let e=a.value;return e===null?0:Math.max(e.getBound(i.value)-1,0)}),s=e=>{let t=a.value;return t===null?0:t.sum(e)},c=Q(()=>{let t=a.value;return t===null?0:Math.min(t.getBound(i.value+r.value)+1,e.value.length-1)});return h(Fe,{startIndexRef:o,endIndexRef:c,columnsRef:e,renderColRef:t,renderItemWithColsRef:n,getLeft:s}),{listWidthRef:r,scrollLeftRef:i}}var Le=r({name:`VirtualListRow`,props:{index:{type:Number,required:!0},item:{type:Object,required:!0}},setup(){let{startIndexRef:e,endIndexRef:t,columnsRef:n,getLeft:r,renderColRef:i,renderItemWithColsRef:a}=f(Fe);return{startIndex:e,endIndex:t,columns:n,renderCol:i,renderItemWithCols:a,getLeft:r}},render(){let{startIndex:e,endIndex:t,columns:n,renderCol:r,renderItemWithCols:i,getLeft:a,item:o}=this;if(i!=null)return i({itemIndex:this.index,startColIndex:e,endColIndex:t,allColumns:n,item:o,getLeft:a});if(r!=null){let i=[];for(let s=e;s<=t;++s){let e=n[s];i.push(r({column:e,left:a(s),item:o}))}return i}return null}}),Re=z(`.v-vl`,{maxHeight:`inherit`,height:`100%`,overflow:`auto`,minWidth:`1px`},[z(`&:not(.v-vl--show-scrollbar)`,{scrollbarWidth:`none`},[z(`&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb`,{width:0,height:0,display:`none`})])]),ze=r({name:`VirtualList`,inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},columns:{type:Array,default:()=>[]},renderCol:Function,renderItemWithCols:Function,items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:`div`},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:`key`},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(e){let t=ne();Re.mount({id:`vueuc/virtual-list`,head:!0,anchorMetaName:U,ssr:t}),s(()=>{let{defaultScrollIndex:t,defaultScrollKey:n}=e;t==null?n!=null&&b({key:n}):b({index:t})});let n=!1,r=!1;a(()=>{if(n=!1,!r){r=!0;return}b({top:_.value,left:l.value})}),i(()=>{n=!0,r||=!0});let o=Q(()=>{if(e.renderCol==null&&e.renderItemWithCols==null||e.columns.length===0)return;let t=0;return e.columns.forEach(e=>{t+=e.width}),t}),c=u(()=>{let t=new Map,{keyField:n}=e;return e.items.forEach((e,r)=>{t.set(e[n],r)}),t}),{scrollLeftRef:l,listWidthRef:d}=Ie({columnsRef:P(e,`columns`),renderColRef:P(e,`renderCol`),renderItemWithColsRef:P(e,`renderItemWithCols`)}),f=N(null),p=N(void 0),m=new Map,h=u(()=>{let{items:t,itemSize:n,keyField:r}=e,i=new Ae(t.length,n);return t.forEach((e,t)=>{let n=e[r],a=m.get(n);a!==void 0&&i.add(t,a)}),i}),g=N(0),_=N(0),v=Q(()=>Math.max(h.value.getBound(_.value-J(e.paddingTop))-1,0)),y=u(()=>{let{value:t}=p;if(t===void 0)return[];let{items:n,itemSize:r}=e,i=v.value,a=Math.min(i+Math.ceil(t/r+1),n.length-1),o=[];for(let e=i;e<=a;++e)o.push(n[e]);return o}),b=(e,t)=>{if(typeof e==`number`){w(e,t,`auto`);return}let{left:n,top:r,index:i,key:a,position:o,behavior:s,debounce:l=!0}=e;if(n!==void 0||r!==void 0)w(n,r,s);else if(i!==void 0)C(i,s,l);else if(a!==void 0){let e=c.value.get(a);e!==void 0&&C(e,s,l)}else o===`bottom`?w(0,2**53-1,s):o===`top`&&w(0,0,s)},x,S=null;function C(t,n,r){let i=f.value;if(i==null)return;let{value:a}=h,o=a.sum(t)+J(e.paddingTop);if(!r)i.scrollTo({left:0,top:o,behavior:n});else{x=t,S!==null&&window.clearTimeout(S),S=window.setTimeout(()=>{x=void 0,S=null},16);let{scrollTop:e,offsetHeight:r}=i;if(o>e){let s=a.get(t);o+s<=e+r||i.scrollTo({left:0,top:o+s-r,behavior:n})}else i.scrollTo({left:0,top:o,behavior:n})}}function w(e,t,n){f.value?.scrollTo({left:e,top:t,behavior:n})}function T(t,r){if(n||e.ignoreItemResize||M(r.target))return;let{value:i}=h,a=c.value.get(t),o=i.get(a),s=r.borderBoxSize?.[0]?.blockSize??r.contentRect.height;if(s===o)return;s-e.itemSize===0?m.delete(t):m.set(t,s-e.itemSize);let l=s-o;if(l===0)return;i.add(a,l);let u=f.value;if(u!=null){if(x===void 0){let e=i.sum(a);u.scrollTop>e&&u.scrollBy(0,l)}else(a<x||a===x&&s+i.sum(a)>u.scrollTop+u.offsetHeight)&&u.scrollBy(0,l);j()}g.value++}let E=!Me(),D=!1;function O(t){var n;(n=e.onScroll)==null||n.call(e,t),(!E||!D)&&j()}function k(t){var n;if((n=e.onWheel)==null||n.call(e,t),E){let e=f.value;if(e!=null){if(t.deltaX===0&&(e.scrollTop===0&&t.deltaY<=0||e.scrollTop+e.offsetHeight>=e.scrollHeight&&t.deltaY>=0))return;t.preventDefault(),e.scrollTop+=t.deltaY/Pe(),e.scrollLeft+=t.deltaX/Pe(),j(),D=!0,xe(()=>{D=!1})}}}function A(t){if(n||M(t.target))return;if(e.renderCol==null&&e.renderItemWithCols==null){if(t.contentRect.height===p.value)return}else if(t.contentRect.height===p.value&&t.contentRect.width===d.value)return;p.value=t.contentRect.height,d.value=t.contentRect.width;let{onResize:r}=e;r!==void 0&&r(t)}function j(){let{value:e}=f;e!=null&&(_.value=e.scrollTop,l.value=e.scrollLeft)}function M(e){let t=e;for(;t!==null;){if(t.style.display===`none`)return!0;t=t.parentElement}return!1}return{listHeight:p,listStyle:{overflow:`auto`},keyToIndex:c,itemsStyle:u(()=>{let{itemResizable:t}=e,n=X(h.value.sum());return g.value,[e.itemsStyle,{boxSizing:`content-box`,width:X(o.value),height:t?``:n,minHeight:t?n:``,paddingTop:X(e.paddingTop),paddingBottom:X(e.paddingBottom)}]}),visibleItemsStyle:u(()=>(g.value,{transform:`translateY(${X(h.value.sum(v.value))})`})),viewportItems:y,listElRef:f,itemsElRef:N(null),scrollTo:b,handleListResize:A,handleListScroll:O,handleListWheel:k,handleItemResize:T}},render(){let{itemResizable:e,keyField:t,keyToIndex:r,visibleItemsTag:i}=this;return c(fe,{onResize:this.handleListResize},{default:()=>{var a;return c(`div`,n(this.$attrs,{class:[`v-vl`,this.showScrollbar&&`v-vl--show-scrollbar`],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:`listElRef`}),[this.items.length===0?(a=this.$slots).empty?.call(a):c(`div`,{ref:`itemsElRef`,class:`v-vl-items`,style:this.itemsStyle},[c(i,Object.assign({class:`v-vl-visible-items`,style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>{let{renderCol:n,renderItemWithCols:i}=this;return this.viewportItems.map(a=>{let o=a[t],s=r.get(o),l=n==null?void 0:c(Le,{index:s,item:a}),u=i==null?void 0:c(Le,{index:s,item:a}),d=this.$slots.default({item:a,renderedCols:l,renderedItemWithCols:u,index:s})[0];return e?c(fe,{key:o,onResize:e=>this.handleItemResize(o,e)},{default:()=>d}):(d.key=o,d)})}})])])}})}});function Be(e){switch(typeof e){case`string`:return e||void 0;case`number`:return String(e);default:return}}function Ve(t,n){n&&(s(()=>{let{value:e}=t;e&&Y.registerHandler(e,n)}),e(t,(e,t)=>{t&&Y.unregisterHandler(t)},{deep:!1}),m(()=>{let{value:e}=t;e&&Y.unregisterHandler(e)}))}var He=r({props:{onFocus:Function,onBlur:Function},setup(e){return()=>(()=>{let t=w(`d16ead82505dc285`);return o(),l(`div`,{style:`width: 0; height: 0`,tabindex:0,onFocus:t[0]||=t=>e.onFocus?.(t),onBlur:t[1]||=t=>e.onBlur?.(t)},null,32)})()}}),Ue=r({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=f(H);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:r,nodeProps:i,tmNode:{rawNode:a}}=this,s=i?.(a),c=t?t(a,!1):$(a[this.labelField],a,!1),u=(o(),l(`div`,n(s,{class:[`${e}-base-select-group-header`,s?.class]}),[L(()=>c)],16));return a.render?a.render({node:u,option:a}):r?r({node:u,option:a,selected:!1}):u}});function We(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var Ge=r({name:`Checkmark`,render(){return(()=>{let e=w(`3c84eac8ae4e1f96`);return e[0]||=E(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},[E(`g`,{fill:`none`},[E(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})])],-1)})()}}),Ke=[`onClick`,`onMouseenter`,`onMousemove`];function qe(e,n){return o(),t(C,{name:`fade-in-scale-up-transition`},{default:()=>e?(o(),t(T,{key:1,clsPrefix:n,class:S(`${n}-base-select-option__check`)},{default:()=>c(Ge)},1032,[`clsPrefix`,`class`])):null},1024)}var Je=r({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:p}=f(H),m=Q(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function h(t){let{tmNode:n}=e;n.disabled||d(t,n)}function g(t){let{tmNode:n}=e;n.disabled||p(t,n)}function _(t){let{tmNode:n}=e,{value:r}=m;n.disabled||r||p(t,n)}return{multiple:r,isGrouped:Q(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:m,isSelected:Q(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:_,handleMouseEnter:g,handleClick:h}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:r,isPending:i,isGrouped:a,showCheckmark:s,nodeProps:c,renderOption:u,renderLabel:d,handleClick:f,handleMouseEnter:p,handleMouseMove:m}=this,h=qe(r,e),g=d?[d(t,r),s&&h]:[$(t[this.labelField],t,r),s&&h],_=c?.(t),v=(o(),l(`div`,n(_,{class:[`${e}-base-select-option`,t.class,_?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:r,[`${e}-base-select-option--grouped`]:a,[`${e}-base-select-option--pending`]:i,[`${e}-base-select-option--show-checkmark`]:s}],style:[_?.style||``,t.style||``],onClick:We([f,_?.onClick]),onMouseenter:We([p,_?.onMouseenter]),onMousemove:We([m,_?.onMousemove])}),[E(`div`,{class:S(`${e}-base-select-option__content`)},[L(()=>g)],2)],16,Ke));return t.render?t.render({node:v,option:t,selected:r}):u?u({node:v,option:t,selected:r}):v}}),Ye=A(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[A(`scrollbar`,`
 max-height: var(--n-height);
 `),A(`virtual-list`,`
 max-height: var(--n-height);
 `),A(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[y(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),A(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),A(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),y(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),y(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),y(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),y(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),A(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),A(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[D(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),_(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),_(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),D(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),D(`pending`,[_(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),D(`selected`,`
 color: var(--n-option-text-color-active);
 `,[_(`&::before`,`
 background-color: var(--n-option-color-active);
 `),D(`pending`,[_(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),D(`disabled`,`
 cursor: not-allowed;
 `,[F(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),D(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),y(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[be({enterScale:`0.5`})])])]),Xe=[`tabindex`,`onFocusin`,`onFocusout`,`onKeyup`,`onKeydown`,`onMousedown`,`onMouseenter`,`onMouseleave`],Ze=r({name:`InternalSelectMenu`,props:{...k.props,clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function},setup(t){let{mergedClsPrefixRef:n,mergedRtlRef:r,mergedComponentPropsRef:i}=te(t),a=M(`InternalSelectMenu`,r,n),o=k(`InternalSelectMenu`,`-internal-select-menu`,Ye,Ee,t,P(t,`clsPrefix`)),c=N(null),l=N(null),d=N(null),f=u(()=>t.treeMate.getFlattenedNodes()),g=u(()=>se(f.value)),_=N(null);function v(){let{treeMate:e}=t,n=null,{value:r}=t;r===null?n=e.getFirstAvailableNode():(n=t.multiple?e.getNode((r||[])[(r||[]).length-1]):e.getNode(r),(!n||n.disabled)&&(n=e.getFirstAvailableNode())),z(n||null)}function y(){let{value:e}=_;e&&!t.treeMate.getNode(e.key)&&(_.value=null)}let x;e(()=>t.show,n=>{n?x=e(()=>t.treeMate,()=>{t.resetMenuOnOptionsChange?(t.autoPending?v():y(),p(B)):y()},{immediate:!0}):x?.()},{immediate:!0}),m(()=>{x?.()});let S=u(()=>J(o.value.self[b(`optionHeight`,t.size)])),C=u(()=>_e(o.value.self[b(`padding`,t.size)])),w=u(()=>t.multiple&&Array.isArray(t.value)?new Set(t.value):new Set),T=u(()=>{let e=f.value;return e&&e.length===0}),E=u(()=>i?.value?.Select?.renderEmpty);function D(e){let{onToggle:n}=t;n&&n(e)}function O(e){let{onScroll:n}=t;n&&n(e)}function A(e){d.value?.sync(),O(e)}function j(){d.value?.sync()}function F(){let{value:e}=_;return e||null}function I(e,t){t.disabled||z(t,!1)}function ne(e,t){t.disabled||D(t)}function L(e){V(e,`action`)||t.onKeyup?.(e)}function re(e){V(e,`action`)||t.onKeydown?.(e)}function ie(e){t.onMousedown?.(e),!t.focusable&&e.preventDefault()}function R(){let{value:e}=_;e&&z(e.getNext({loop:!0}),!0)}function ae(){let{value:e}=_;e&&z(e.getPrev({loop:!0}),!0)}function z(e,t=!1){_.value=e,t&&B()}function B(){let e=_.value;if(!e)return;let n=g.value(e.key);n!==null&&(t.virtualScroll?l.value?.scrollTo({index:n}):d.value?.scrollTo({index:n,elSize:S.value}))}function oe(e){c.value?.contains(e.target)&&t.onFocus?.(e)}function U(e){c.value?.contains(e.relatedTarget)||t.onBlur?.(e)}h(H,{handleOptionMouseEnter:I,handleOptionClick:ne,valueSetRef:w,pendingTmNodeRef:_,nodePropsRef:P(t,`nodeProps`),showCheckmarkRef:P(t,`showCheckmark`),multipleRef:P(t,`multiple`),valueRef:P(t,`value`),renderLabelRef:P(t,`renderLabel`),renderOptionRef:P(t,`renderOption`),labelFieldRef:P(t,`labelField`),valueFieldRef:P(t,`valueField`)}),h(W,c),s(()=>{let{value:e}=d;e&&e.sync()});let G=u(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:a,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:x,[b(`optionFontSize`,e)]:S,[b(`optionHeight`,e)]:C,[b(`optionPadding`,e)]:w}}=o.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":a,"--n-option-font-size":S,"--n-group-header-text-color":s,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":x,"--n-option-height":C,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":w,"--n-option-padding-left":_e(w,`left`),"--n-option-padding-right":_e(w,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:ce}=t,K=ce?ee(`internal-select-menu`,u(()=>t.size[0]),G,t):void 0,le={selfRef:c,next:R,prev:ae,getPendingTmNode:F};return Ve(c,t.onResize),{mergedTheme:o,mergedClsPrefix:n,rtlEnabled:a,virtualListRef:l,scrollbarRef:d,itemSize:S,padding:C,flattenedNodes:f,empty:T,mergedRenderEmpty:E,virtualListContainer(){let{value:e}=l;return e?.listElRef},virtualListContent(){let{value:e}=l;return e?.itemsElRef},doScroll:O,handleFocusin:oe,handleFocusout:U,handleKeyUp:L,handleKeyDown:re,handleMouseDown:ie,handleVirtualListResize:j,handleVirtualListScroll:A,cssVars:ce?void 0:G,themeClass:K?.themeClass,onRender:K?.onRender,...le}},render(){let{$slots:e,virtualScroll:r,clsPrefix:i,mergedTheme:a,themeClass:s,onRender:c}=this;return c?.(),o(),l(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:S([`${i}-base-select-menu`,`${i}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${i}-base-select-menu--rtl`,s,this.multiple&&`${i}-base-select-menu--multiple`]),style:j(this.cssVars),onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},[L(()=>me(e.header,e=>e&&(o(),l(`div`,{class:S(`${i}-base-select-menu__header`),"data-header":!0,key:`header`},[L(()=>e)],2)))),this.loading?(o(),l(`div`,{key:0,class:S(`${i}-base-select-menu__loading`)},[(o(),t(Oe,{clsPrefix:i,strokeWidth:20},null,8,[`clsPrefix`]))],2)):(o(),l(I,{key:1},[this.empty?(o(),l(`div`,{key:1,class:S(`${i}-base-select-menu__empty`),"data-empty":!0},[L(()=>de(e.empty,()=>[this.mergedRenderEmpty?.()||(o(),t(ie,{theme:a.peers.Empty,themeOverrides:a.peerOverrides.Empty,size:this.size},null,8,[`theme`,`themeOverrides`,`size`]))]))],2)):(o(),t(pe,n({key:0,ref:`scrollbarRef`,theme:a.peers.Scrollbar,themeOverrides:a.peerOverrides.Scrollbar,scrollable:this.scrollable,container:r?this.virtualListContainer:void 0,content:r?this.virtualListContent:void 0,onScroll:r?void 0:this.doScroll},this.scrollbarProps),{default:()=>r?(o(),t(ze,{key:1,ref:`virtualListRef`,class:S(`${i}-virtual-list`),items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?(o(),t(Ue,{key:e.key,clsPrefix:i,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):e.ignored?null:(o(),t(Je,{clsPrefix:i,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))},1032,[`class`,`items`,`itemSize`,`paddingTop`,`paddingBottom`,`onResize`,`onScroll`])):(o(),l(`div`,{key:4,class:S(`${i}-base-select-menu-option-wrapper`),style:j({paddingTop:this.padding.top,paddingBottom:this.padding.bottom})},[L(()=>this.flattenedNodes.map(e=>e.isGroup?(o(),t(Ue,{key:e.key,clsPrefix:i,tmNode:e},null,8,[`clsPrefix`,`tmNode`])):(o(),t(Je,{clsPrefix:i,key:e.key,tmNode:e},null,8,[`clsPrefix`,`tmNode`]))))],6))},1040,[`theme`,`themeOverrides`,`scrollable`,`container`,`content`,`onScroll`]))],64)),L(()=>me(e.action,e=>e&&[(o(),l(`div`,{class:S(`${i}-base-select-menu__action`),"data-action":!0,key:`action`},[L(()=>e)],2)),(o(),t(He,{onFocus:this.onTabOut,key:`focus-detector`},null,8,[`onFocus`]))]))],46,Xe)}});function Qe(e){return e.type===`group`}function $e(e){return e.type===`ignored`}function et(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function tt(e,t){return{getIsGroup:Qe,getIgnored:$e,getKey(t){return Qe(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function nt(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(Qe(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if($e(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function rt(e,t,n){let r=new Map;return e.forEach(e=>{Qe(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var it=_([A(`base-selection`,`
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
 `,[A(`base-loading`,`
 color: var(--n-loading-color);
 `),A(`base-selection-tags`,`min-height: var(--n-height);`),y(`border, state-border`,`
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
 `),y(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),A(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[y(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),A(`base-selection-overlay`,`
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
 `,[y(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),A(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[y(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),A(`base-selection-tags`,`
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
 `),A(`base-selection-label`,`
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
 `,[A(`base-selection-input`,`
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
 `,[y(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),y(`render-label`,`
 color: var(--n-text-color);
 `)]),F(`disabled`,[_(`&:hover`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),D(`focus`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),D(`active`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),A(`base-selection-label`,`background-color: var(--n-color-active);`),A(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),D(`disabled`,`cursor: not-allowed;`,[y(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),A(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[A(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),y(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),A(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),A(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),A(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[y(`input`,`
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
 `),y(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>D(`${e}-status`,[y(`state-border`,`border: var(--n-border-${e});`),F(`disabled`,[_(`&:hover`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),D(`active`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),A(`base-selection-label`,`background-color: var(--n-color-active-${e});`),A(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),D(`focus`,[y(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),A(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),A(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[_(`&:last-child`,`padding-right: 0;`),A(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[y(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),at=[`disabled`,`value`,`autofocus`,`onBlur`,`onFocus`,`onKeydown`,`onInput`,`onCompositionstart`,`onCompositionend`],ot=[`tabindex`],st=[`title`],ct=[`value`,`readonly`,`disabled`,`autofocus`,`onFocus`,`onBlur`,`onInput`,`onCompositionstart`,`onCompositionend`],lt=[`tabindex`],ut=[`onClick`,`onMouseenter`,`onMouseleave`,`onKeydown`,`onFocusin`,`onFocusout`,`onMousedown`],dt=r({name:`InternalSelection`,props:{...k.props,clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function},setup(t){let{mergedClsPrefixRef:n,mergedRtlRef:r}=te(t),i=M(`InternalSelection`,r,n),a=N(null),o=N(null),c=N(null),l=N(null),d=N(null),f=N(null),m=N(null),h=N(null),g=N(null),_=N(null),v=N(!1),y=N(!1),S=N(!1),C=k(`InternalSelection`,`-internal-selection`,it,Te,t,P(t,`clsPrefix`)),w=u(()=>t.clearable&&!t.disabled&&(S.value||t.active)),T=u(()=>t.selectedOption?t.renderTag?t.renderTag({option:t.selectedOption,handleClose:()=>{}}):t.renderLabel?t.renderLabel(t.selectedOption,!0):$(t.selectedOption[t.labelField],t.selectedOption,!0):t.placeholder),E=u(()=>{let e=t.selectedOption;if(e)return e[t.labelField]}),D=u(()=>t.multiple?!!(Array.isArray(t.selectedOptions)&&t.selectedOptions.length):t.selectedOption!==null);function O(){let{value:e}=a;if(e){let{value:n}=o;n&&(n.style.width=`${e.offsetWidth}px`,t.maxTagCount!==`responsive`&&g.value?.sync({showAllItemsBeforeCalculate:!1}))}}function A(){let{value:e}=_;e&&(e.style.display=`none`)}function j(){let{value:e}=_;e&&(e.style.display=`inline-block`)}e(P(t,`active`),e=>{e||A()}),e(P(t,`pattern`),()=>{t.multiple&&p(O)});function F(e){let{onFocus:n}=t;n&&n(e)}function I(e){let{onBlur:n}=t;n&&n(e)}function ne(e){let{onDeleteOption:n}=t;n&&n(e)}function L(e){let{onClear:n}=t;n&&n(e)}function re(e){let{onPatternInput:n}=t;n&&n(e)}function ie(e){(!e.relatedTarget||!c.value?.contains(e.relatedTarget))&&F(e)}function R(e){c.value?.contains(e.relatedTarget)||I(e)}function ae(e){L(e)}function z(){S.value=!0}function B(){S.value=!1}function oe(e){t.active&&t.filterable&&e.target!==o.value&&e.preventDefault()}function V(e){ne(e)}let H=N(!1);function U(e){if(e.key===`Backspace`&&!H.value&&!t.pattern.length){let{selectedOptions:e}=t;e?.length&&V(e[e.length-1])}}let W=null;function se(e){let{value:n}=a;n&&(n.textContent=e.target.value,O()),t.ignoreComposition&&H.value?W=e:re(e)}function G(){H.value=!0}function ce(){H.value=!1,t.ignoreComposition&&re(W),W=null}function K(e){y.value=!0,t.onPatternFocus?.(e)}function le(e){y.value=!1,t.onPatternBlur?.(e)}function ue(){if(t.filterable)y.value=!1,f.value?.blur(),o.value?.blur();else if(t.multiple){let{value:e}=l;e?.blur()}else{let{value:e}=d;e?.blur()}}function q(){t.filterable?(y.value=!1,f.value?.focus()):t.multiple?l.value?.focus():d.value?.focus()}function J(){let{value:e}=o;e&&(j(),e.focus())}function Y(){let{value:e}=o;e&&e.blur()}function X(e){let{value:t}=m;t&&t.setTextContent(`+${e}`)}function de(){let{value:e}=h;return e}function fe(){return o.value}let Z=null;function pe(){Z!==null&&window.clearTimeout(Z)}function me(){t.active||(pe(),Z=window.setTimeout(()=>{D.value&&(v.value=!0)},100))}function he(){pe()}function ge(e){e||(pe(),v.value=!1)}e(D,e=>{e||(v.value=!1)}),s(()=>{x(()=>{let e=f.value;e&&(t.disabled?e.removeAttribute(`tabindex`):e.tabIndex=y.value?-1:0)})}),Ve(c,t.onResize);let{inlineThemeDisabled:ve}=t,ye=u(()=>{let{size:e}=t,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:y,borderHover:x,borderActive:S,arrowColor:w,arrowColorDisabled:T,loadingColor:E,colorActiveWarning:D,boxShadowFocusWarning:O,boxShadowActiveWarning:k,boxShadowHoverWarning:A,borderWarning:j,borderFocusWarning:M,borderHoverWarning:ee,borderActiveWarning:N,colorActiveError:P,boxShadowFocusError:F,boxShadowActiveError:te,boxShadowHoverError:I,borderError:ne,borderFocusError:L,borderHoverError:re,borderActiveError:ie,clearColor:R,clearColorHover:ae,clearColorPressed:z,clearSize:B,arrowSize:oe,[b(`height`,e)]:V,[b(`fontSize`,e)]:H}}=C.value,U=_e(c),W=_e(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":y,"--n-border-hover":x,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":H,"--n-height":V,"--n-padding-single-top":U.top,"--n-padding-multiple-top":W.top,"--n-padding-single-right":U.right,"--n-padding-multiple-right":W.right,"--n-padding-single-left":U.left,"--n-padding-multiple-left":W.left,"--n-padding-single-bottom":U.bottom,"--n-padding-multiple-bottom":W.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":w,"--n-arrow-color-disabled":T,"--n-loading-color":E,"--n-color-active-warning":D,"--n-box-shadow-focus-warning":O,"--n-box-shadow-active-warning":k,"--n-box-shadow-hover-warning":A,"--n-border-warning":j,"--n-border-focus-warning":M,"--n-border-hover-warning":ee,"--n-border-active-warning":N,"--n-color-active-error":P,"--n-box-shadow-focus-error":F,"--n-box-shadow-active-error":te,"--n-box-shadow-hover-error":I,"--n-border-error":ne,"--n-border-focus-error":L,"--n-border-hover-error":re,"--n-border-active-error":ie,"--n-clear-size":B,"--n-clear-color":R,"--n-clear-color-hover":ae,"--n-clear-color-pressed":z,"--n-arrow-size":oe,"--n-font-weight":r}}),be=ve?ee(`internal-selection`,u(()=>t.size[0]),ye,t):void 0;return{mergedTheme:C,mergedClearable:w,mergedClsPrefix:n,rtlEnabled:i,patternInputFocused:y,filterablePlaceholder:T,label:E,selected:D,showTagsPanel:v,isComposing:H,counterRef:m,counterWrapperRef:h,patternInputMirrorRef:a,patternInputRef:o,selfRef:c,multipleElRef:l,singleElRef:d,patternInputWrapperRef:f,overflowRef:g,inputTagElRef:_,handleMouseDown:oe,handleFocusin:ie,handleClear:ae,handleMouseEnter:z,handleMouseLeave:B,handleDeleteOption:V,handlePatternKeyDown:U,handlePatternInputInput:se,handlePatternInputBlur:le,handlePatternInputFocus:K,handleMouseEnterCounter:me,handleMouseLeaveCounter:he,handleFocusout:R,handleCompositionEnd:ce,handleCompositionStart:G,onPopoverUpdateShow:ge,focus:q,focusInput:J,blur:ue,blurInput:Y,updateCounter:X,getCounter:de,getTail:fe,renderLabel:t.renderLabel,cssVars:ve?void 0:ye,themeClass:be?.themeClass,onRender:be?.onRender}},render(){let{status:e,multiple:r,size:i,disabled:a,filterable:s,maxTagCount:c,bordered:u,clsPrefix:d,ellipsisTagPopoverProps:f,onRender:p,renderTag:m,renderLabel:h}=this;p?.();let g=c===`responsive`,_=typeof c==`number`,v=g||_,y=(o(),t(he,null,{default:()=>(o(),t(R,{clsPrefix:d,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>this.$slots.arrow?.()},1032,[`clsPrefix`,`loading`,`showArrow`,`showClear`,`onClear`]))},1024)),b;if(r){let{labelField:e}=this,r=n=>(o(),l(`div`,{class:S(`${d}-base-selection-tag-wrapper`),key:n.value},[m?(o(),l(I,{key:0},[L(()=>m({option:n,handleClose:()=>{this.handleDeleteOption(n)}}))],64)):(o(),t(Ce,{key:1,size:i,closable:!n.disabled,disabled:a,onClose:()=>{this.handleDeleteOption(n)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>h?h(n,!0):$(n[e],n,!0)},1032,[`size`,`closable`,`disabled`,`onClose`]))],2)),u=()=>(_?this.selectedOptions.slice(0,c):this.selectedOptions).map(r),p=s?(o(),l(`div`,{class:S(`${d}-base-selection-input-tag`),ref:`inputTagElRef`,key:`__input-tag__`},[E(`input`,n(this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:a,value:this.pattern,autofocus:this.autofocus,class:`${d}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,at),E(`span`,{ref:`patternInputMirrorRef`,class:S(`${d}-base-selection-input-tag__mirror`)},[L(()=>this.pattern)],2)],2)):null,x=g?()=>(o(),l(`div`,{class:S(`${d}-base-selection-tag-wrapper`),ref:`counterWrapperRef`},[(o(),t(Ce,{size:i,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:a},null,8,[`size`,`onMouseenter`,`onMouseleave`,`disabled`]))],2)):void 0,C;if(_){let e=this.selectedOptions.length-c;e>0&&(C=(n=>(o(),l(`div`,{class:S(`${d}-base-selection-tag-wrapper`),key:`__counter__`},[(o(),t(Ce,{size:i,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:a},{default:()=>`+${e}`},1032,[`size`,`onMouseenter`,`disabled`]))],2)))(C))}let w=g?s?(o(),t(G,{key:3,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:u,counter:x,tail:()=>p},1032,[`updateCounter`,`getCounter`,`getTail`])):(o(),t(G,{key:4,ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:u,counter:x},1032,[`updateCounter`,`getCounter`])):_&&C?u().concat(C):u(),T=v?()=>(o(),l(`div`,{class:S(`${d}-base-selection-popover`)},[g?(o(),l(I,{key:0},[L(()=>u())],64)):(o(),l(I,{key:1},[L(()=>this.selectedOptions.map(r))],64))],2)):void 0,D=v?{show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover,...f}:null,O=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?(o(),l(`div`,{key:5,class:S(`${d}-base-selection-placeholder ${d}-base-selection-overlay`)},[E(`div`,{class:S(`${d}-base-selection-placeholder__inner`)},[L(()=>this.placeholder)],2)],2)):null,k=s?(o(),l(`div`,{key:6,ref:`patternInputWrapperRef`,class:S(`${d}-base-selection-tags`)},[L(()=>w),g?L(()=>null):(o(),l(I,{key:1},[L(()=>p)],64)),L(()=>y)],2)):(o(),l(`div`,{key:7,ref:`multipleElRef`,class:S(`${d}-base-selection-tags`),tabindex:a?void 0:0},[L(()=>w),L(()=>y)],10,ot));b=(e=>(o(),l(I,{key:8},[v?(o(),t(K,n({key:0},D,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>k,default:T},1040)):(o(),l(I,{key:1},[L(()=>k)],64)),L(()=>O)],64)))(b)}else if(s){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,r=!this.active&&this.selected;b=(e=>(o(),l(`div`,{key:9,ref:`patternInputWrapperRef`,class:S(`${d}-base-selection-label`),title:this.patternInputFocused?void 0:Be(this.label)},[E(`input`,n(this.inputProps,{ref:`patternInputRef`,class:`${d}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:a,disabled:a,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd}),null,16,ct),r?(o(),l(`div`,{class:S(`${d}-base-selection-label__render-label ${d}-base-selection-overlay`),key:`input`},[E(`div`,{class:S(`${d}-base-selection-overlay__wrapper`)},[m?(o(),l(I,{key:0},[L(()=>m({option:this.selectedOption,handleClose:()=>{}}))],64)):(o(),l(I,{key:1},[h?(o(),l(I,{key:0},[L(()=>h(this.selectedOption,!0))],64)):(o(),l(I,{key:1},[L(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],2)):L(()=>null),t?(o(),l(`div`,{class:S(`${d}-base-selection-placeholder ${d}-base-selection-overlay`),key:`placeholder`},[E(`div`,{class:S(`${d}-base-selection-overlay__wrapper`)},[L(()=>this.filterablePlaceholder)],2)],2)):L(()=>null),L(()=>y)],10,st)))(b)}else b=(e=>(o(),l(`div`,{key:10,ref:`singleElRef`,class:S(`${d}-base-selection-label`),tabindex:this.disabled?void 0:0},[this.label===void 0?(o(),l(`div`,{class:S(`${d}-base-selection-placeholder ${d}-base-selection-overlay`),key:`placeholder`},[E(`div`,{class:S(`${d}-base-selection-placeholder__inner`)},[L(()=>this.placeholder)],2)],2)):(o(),l(`div`,{class:S(`${d}-base-selection-input`),title:Be(this.label),key:`input`},[E(`div`,{class:S(`${d}-base-selection-input__content`)},[m?(o(),l(I,{key:0},[L(()=>m({option:this.selectedOption,handleClose:()=>{}}))],64)):(o(),l(I,{key:1},[h?(o(),l(I,{key:0},[L(()=>h(this.selectedOption,!0))],64)):(o(),l(I,{key:1},[L(()=>$(this.label,this.selectedOption,!0))],64))],64))],2)],10,[`title`])),L(()=>y)],10,lt)))(b);return o(),l(`div`,{ref:`selfRef`,class:S([`${d}-base-selection`,this.rtlEnabled&&`${d}-base-selection--rtl`,this.themeClass,e&&`${d}-base-selection--${e}-status`,{[`${d}-base-selection--active`]:this.active,[`${d}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${d}-base-selection--disabled`]:this.disabled,[`${d}-base-selection--multiple`]:this.multiple,[`${d}-base-selection--focus`]:this.focused}]),style:j(this.cssVars),onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},[L(()=>b),u?(o(),l(`div`,{key:0,class:S(`${d}-base-selection__border`)},null,2)):L(()=>null),u?(o(),l(`div`,{key:2,class:S(`${d}-base-selection__state-border`)},null,2)):L(()=>null)],46,ut)}}),ft=_([A(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),A(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[be({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),pt={...k.props,to:oe.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array},mt=r({name:`Select`,props:pt,slots:Object,setup(t){let{mergedClsPrefixRef:n,mergedBorderedRef:r,namespaceRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:o}=te(t),s=k(`Select`,`-select`,ft,De,t,n),c=N(t.defaultValue),l=P(t,`value`),d=ye(l,c),f=N(!1),p=N(``),m=ce(t,[`items`,`options`]),h=N([]),g=N([]),_=u(()=>g.value.concat(h.value).concat(m.value)),v=u(()=>{let{filter:e}=t;if(e)return e;let{labelField:n,valueField:r}=t;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return et(e,i);let a=t[r];return typeof a==`string`?et(e,a):typeof a==`number`&&et(e,String(a))}}),y=u(()=>{if(t.remote)return m.value;{let{value:e}=_,{value:n}=p;return!n.length||!t.filterable?e:nt(e,v.value,n,t.childrenField)}}),b=u(()=>{let{valueField:e,childrenField:n}=t,r=tt(e,n);return ue(y.value,r)}),x=u(()=>rt(_.value,t.valueField,t.childrenField)),S=N(!1),C=ye(P(t,`show`),S),w=N(null),T=N(null),E=N(null),{localeRef:D}=ae(`Select`),O=u(()=>t.placeholder??D.value.placeholder),A=[],j=N(new Map),M=u(()=>{let{fallbackOption:e}=t;if(e===void 0){let{labelField:e,valueField:n}=t;return t=>({[e]:String(t),[n]:t})}return e===!1?!1:t=>Object.assign(e(t),{value:t})});function F(e){let n=t.remote,{value:r}=j,{value:i}=x,{value:a}=M,o=[];return e.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let I=u(()=>{if(t.multiple){let{value:e}=d;return Array.isArray(e)?F(e):[]}return null}),ne=u(()=>{let{value:e}=d;return!t.multiple&&!Array.isArray(e)?e===null?null:F([e])[0]||null:null}),L=we(t,{mergedSize:e=>{let{size:n}=t;if(n)return n;let{mergedSize:r}=e||{};return r?.value?r.value:o?.value?.Select?.size||`medium`}}),{mergedSizeRef:ie,mergedDisabledRef:R,mergedStatusRef:z}=L;function B(e,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=t,{nTriggerFormChange:o,nTriggerFormInput:s}=L;r&&Z(r,e,n),a&&Z(a,e,n),i&&Z(i,e,n),c.value=e,o(),s()}function H(e){let{onBlur:n}=t,{nTriggerFormBlur:r}=L;n&&Z(n,e),r()}function U(){let{onClear:e}=t;e&&Z(e)}function W(e){let{onFocus:n,showOnFocus:r}=t,{nTriggerFormFocus:i}=L;n&&Z(n,e),i(),r&&q()}function se(e){let{onSearch:n}=t;n&&Z(n,e)}function G(e){let{onScroll:n}=t;n&&Z(n,e)}function K(){let{remote:e,multiple:n}=t;if(e){let{value:e}=j;if(n){let{valueField:n}=t;I.value?.forEach(t=>{e.set(t[n],t)})}else{let n=ne.value;n&&e.set(n[t.valueField],n)}}}function le(e){let{onUpdateShow:n,"onUpdate:show":r}=t;n&&Z(n,e),r&&Z(r,e),S.value=e}function q(){R.value||(le(!0),S.value=!0,t.filterable&&Ne())}function J(){le(!1)}function Y(){p.value=``,g.value=A}let X=N(!1);function de(){t.filterable&&(X.value=!0)}function fe(){t.filterable&&(X.value=!1,C.value||Y())}function pe(){R.value||(C.value?t.filterable?Ne():J():q())}function me(e){E.value?.selfRef?.contains(e.relatedTarget)||(f.value=!1,H(e),J())}function he(e){W(e),f.value=!0}function _e(){f.value=!0}function ve(e){w.value?.$el.contains(e.relatedTarget)||(f.value=!1,H(e),J())}function be(){w.value?.focus(),J()}function xe(e){C.value&&(w.value?.$el.contains(ge(e))||J())}function Q(e){if(!Array.isArray(e))return[];if(M.value)return Array.from(e);{let{remote:n}=t,{value:r}=x;if(n){let{value:t}=j;return e.filter(e=>r.has(e)||t.has(e))}return e.filter(e=>r.has(e))}}function Ce(e){Te(e.rawNode)}function Te(e){if(R.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=t;if(n&&!r){let{value:e}=g,t=e[0]||null;if(t){let e=h.value;e.length?e.push(t):h.value=[t],g.value=A}}if(r&&j.value.set(e[a],e),t.multiple){let t=Q(d.value),o=t.findIndex(t=>t===e[a]);if(~o){if(t.splice(o,1),n&&!r){let t=$(e[a]);~t&&(h.value.splice(t,1),i&&(p.value=``))}}else t.push(e[a]),i&&(p.value=``);B(t,F(t))}else{if(n&&!r){let t=$(e[a]);~t?h.value=[h.value[t]]:h.value=A}Me(),J(),B(e[a],e)}}function $(e){return h.value.findIndex(n=>n[t.valueField]===e)}function Ee(e){C.value||q();let{value:n}=e.target;p.value=n;let{tag:r,remote:i}=t;if(se(n),r&&!i){if(!n){g.value=A;return}let{onCreate:e}=t,r=e?e(n):{[t.labelField]:n,[t.valueField]:n},{valueField:i,labelField:a}=t;m.value.some(e=>e[i]===r[i]||e[a]===r[a])||h.value.some(e=>e[i]===r[i]||e[a]===r[a])?g.value=A:g.value=[r]}}function Oe(e){e.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=t;!n&&t.filterable&&J(),r&&!i&&a&&(h.value=A),U(),n?B([],[]):B(null,null)}function ke(e){!V(e,`action`)&&!V(e,`empty`)&&!V(e,`header`)&&e.preventDefault()}function Ae(e){G(e)}function je(e){if(!t.keyboard){e.preventDefault();return}switch(e.key){case` `:if(t.filterable)break;e.preventDefault();case`Enter`:if(!w.value?.isComposing){if(C.value){let e=E.value?.getPendingTmNode();e?Ce(e):t.filterable||(J(),Me())}else if(q(),t.tag&&X.value){let e=g.value[0];if(e){let n=e[t.valueField],{value:r}=d;t.multiple&&Array.isArray(r)&&r.includes(n)||Te(e)}}}e.preventDefault();break;case`ArrowUp`:if(e.preventDefault(),t.loading)return;C.value&&E.value?.prev();break;case`ArrowDown`:if(e.preventDefault(),t.loading)return;C.value?E.value?.next():q();break;case`Escape`:C.value&&(re(e),J()),w.value?.focus()}}function Me(){w.value?.focus()}function Ne(){w.value?.focusInput()}function Pe(){C.value&&T.value?.syncPosition()}K(),e(P(t,`options`),K);let Fe={focus:()=>{w.value?.focus()},focusInput:()=>{w.value?.focusInput()},blur:()=>{w.value?.blur()},blurInput:()=>{w.value?.blurInput()}},Ie=u(()=>{let{self:{menuBoxShadow:e}}=s.value;return{"--n-menu-box-shadow":e}}),Le=a?ee(`select`,void 0,Ie,t):void 0;return{...Fe,mergedStatus:z,mergedClsPrefix:n,mergedBordered:r,namespace:i,treeMate:b,isMounted:Se(),triggerRef:w,menuRef:E,pattern:p,uncontrolledShow:S,mergedShow:C,adjustedTo:oe(t),uncontrolledValue:c,mergedValue:d,followerRef:T,localizedPlaceholder:O,selectedOption:ne,selectedOptions:I,mergedSize:ie,mergedDisabled:R,focused:f,activeWithoutMenuOpen:X,inlineThemeDisabled:a,onTriggerInputFocus:de,onTriggerInputBlur:fe,handleTriggerOrMenuResize:Pe,handleMenuFocus:_e,handleMenuBlur:ve,handleMenuTabOut:be,handleTriggerClick:pe,handleToggle:Ce,handleDeleteOption:Te,handlePatternInput:Ee,handleClear:Oe,handleTriggerBlur:me,handleTriggerFocus:he,handleKeydown:je,handleMenuAfterLeave:Y,handleMenuClickOutside:xe,handleMenuScroll:Ae,handleMenuKeydown:je,handleMenuMousedown:ke,mergedTheme:s,cssVars:a?void 0:Ie,themeClass:Le?.themeClass,onRender:Le?.onRender}},render(){return o(),l(`div`,{class:S(`${this.mergedClsPrefix}-select`)},[d(B,null,{_:1,default:g(()=>[(o(),t(q,null,{_:1,default:g(()=>(o(),t(dt,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{_:1,arrow:g(()=>[this.$slots.arrow?.()])},8,`inlineThemeDisabled.status.inputProps.clsPrefix.showArrow.maxTagCount.ellipsisTagPopoverProps.bordered.active.pattern.placeholder.selectedOption.selectedOptions.multiple.renderTag.renderLabel.filterable.clearable.disabled.size.theme.labelField.valueField.themeOverrides.loading.focused.onClick.onDeleteOption.onPatternInput.onClear.onBlur.onFocus.onKeydown.onPatternBlur.onPatternFocus.onResize.ignoreComposition`.split(`.`))))})),(o(),t(le,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===oe.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{_:1,default:g(()=>(o(),t(C,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{_:1,default:g(()=>this.mergedShow||this.displayDirective===`show`?(this.onRender?.(),O((o(),t(Ze,n(this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{_:1,empty:g(()=>[this.$slots.empty?.()]),header:g(()=>[this.$slots.header?.()]),action:g(()=>[this.$slots.action?.()])},16,`onResize.inlineThemeDisabled.virtualScroll.class.clsPrefix.labelField.valueField.nodeProps.theme.themeOverrides.treeMate.multiple.size.renderOption.renderLabel.value.style.onToggle.onScroll.onFocus.onBlur.onKeydown.onTabOut.onMousedown.show.showCheckmark.resetMenuOnOptionsChange.scrollbarProps`.split(`.`))),this.displayDirective===`show`?[[v,this.mergedShow],[ve,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[ve,this.handleMenuClickOutside,void 0,{capture:!0}]])):null)},8,[`appear`,`onAfterLeave`])))},8,[`show`,`to`,`teleportDisabled`,`containerClass`,`width`,`placement`]))])})],2)}});export{ze as a,We as i,tt as n,Ze as r,mt as t};