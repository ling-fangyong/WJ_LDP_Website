<!-- 数据分析校正结果界面-->
<template>
    <div class="Analysis" v-loading = "loading" element-loading-text="加载中，请稍等...">
        <el-card class="box-card" v-for="(item,index) in QuesAndOp" :key="(index+1).toString()" style="margin 10px;">
            <div slot="header" class="clearfix">
                <div class="quesTitle">
                <span style="color: black; margin-right:3px;"> {{(index+1)+'.'}}  </span>
                {{item.title}}
                <span v-show="item.type==3">&nbsp;&nbsp;&nbsp;&nbsp;(MinValue:&nbsp;{{item.DataMin}}&nbsp;&nbsp;&nbsp;&nbsp;MaxValue:&nbsp;{{item.DataMax}})</span>
                </div>
            </div>
            <template v-if=" item.type == 1 || item.type == 2">
                <el-table :data="item.options" style="width:100%">
                    <el-table-column prop="title" label="选项" ></el-table-column>
                    <el-table-column prop="CalcOp" label="数量"></el-table-column>
                </el-table>
            </template>
            <template v-if=" item.type == 3">
                <el-col :span="12">
                    <span>平均值</span>
                </el-col>
                <el-col :span="12">
                    <template v-if="parseFloat(item.textValue) == NaN">
                        <span>{{item.textValue}}</span>
                    </template>
                    <template v-else>
                        <span>{{parseFloat(item.textValue)}}</span>
                    </template>
                    
                </el-col>
            </template>
        </el-card>
    </div>  
</template>

<script>
import * as API from '../api/answer'

export default ({
    data(){
        return{
            loading:false,
            WjId:0,
            QuesAndOp:[],

        }
    },
    methods: {
        init(Id){
            this.WjId = Id;
            this.QuesAndOp=[];
            this.loading=true;
            var Formdata = new FormData();
            Formdata.append('WjId',this.WjId.toString());
            API.AnalysisData(Formdata).then(res =>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.QuesAndOp = res.data;
                    console.log(this.QuesAndOp);
                    this.loading=false;
                }else{
                    this.$message({
                        type:"error",
                        message:res.msg+'!',
                    })
                }
            }).catch(error=>{
                this.$notify.error({
                    message:error,
                })
            })
        },
    },
})
</script>
