import React, {Component} from "react"


class Preloader extends Component {
    render(){
        return (
        <div class="preloader"  >
                <div class="preloader-inner" style={{ "position": "relative", "top": "unset", "left": "unset", "paddingTop": "45vh" }}>
                <div></div>
                <hr />
            </div>
        </div>
        )
    }
}

export default Preloader;