import { Component } from "react";
import AuthTopNav from "../components/AuthTopNav";
import v1 from "../api/v1";
import ViewPayment from "./ViewPayment";

class PublicCreatorNotVerified extends Component{
    constructor(props) {
        super(props);
        this.state = {
		}
		
		v1.page.track()
		v1.page.event('Creator', 'Not Verified', window.location.href)
        
    }
    
	render() {
		return (<>
			<div class="container custom_width" style={{minHeight: "50vh", marginTop: "15vh"}}>
						<div class="title-block">
							<h2 class="subtitle--about text-black text-center">Sorry this user is not verified</h2>
						</div>
						<p class="about-text text-center text-black ">
							We restrict functionality of MyHustle to verified users only.<br/>
						</p>
						<div class="row">
							<div class="col-md-4 col-sm-4">
								<div class="image-box-wrapper">
									<div class="image-box">
										<a href="/" class="image-box-link">
											<img src="/public/static/images/thumbnail-1.jpg" class="img-responsive image-box-img" alt="" title="imagebox" />
										</a>
									</div>
									<div class="title-block">
										<h6 class="subtitle text-black no-pd-bt">Home Page</h6>
										<p></p>
									</div>
								</div>
							</div>
							<div class="col-md-4 col-sm-4">
									<div class="image-box-wrapper">
										<div class="image-box">
											<a href="/auth/login" class="image-box-link">
												<img src="/public/static/images/thumbnail-2.jpg" class="img-responsive image-box-img" alt="" title="imagebox" />
											</a>
										</div>
										<div class="title-block">
											<h6 class="subtitle text-black no-pd-bt">Login to your account</h6>
										</div>
									</div>
							</div>
							<div class="col-md-4 col-sm-4">
								<div class="image-box-wrapper">
									<div class="image-box">
										<a href="/about" class="image-box-link">
											<img src="/public/static/images/thumbnail-3.jpg" class="img-responsive image-box-img" alt="" title="imagebox" />
										</a>
									</div>
									<div class="title-block">
										<h6 class="subtitle text-black no-pd-bt">About MyHustle</h6>
									</div>
								</div>
							</div>
						</div>
		</div>
		</>)}
}

export default PublicCreatorNotVerified;