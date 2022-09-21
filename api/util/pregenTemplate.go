package util

var pregenTemplate = `
<script>
try{(function(){
window.pregenerated_%s = %s;
})()}catch(e){
	console.log('Failed to render preload with error '+ e)
}finally{ console.log(window.pregenerated_%s)}
</script>
`
