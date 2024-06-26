import * as React from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

export default function AlertDialog() {
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (

	<React.Fragment>
  <Button 
    variant="outlined" 
    style={{ 
      backgroundColor: '#0f172a', 
      textTransform: 'none', 
      fontWeight: '400', 
      fontSize: '16px', 
      color: '#cbd5e1', 
      border: "1px solid #475569",
      marginTop: "18px"
    }} 
    onClick={handleClickOpen}
  >
    How Does This Work?
  </Button>

  <Dialog
    open={open}
    onClose={handleClose}
    aria-labelledby="alert-dialog-title"
    aria-describedby="alert-dialog-description"
    PaperProps={{
		style: { backgroundColor: '#1a1f29', color: '#ffffff',  border: '1px white solid', borderRadius: '17px'}, // Setting the text color to white
    }}
  >
    <DialogTitle id="alert-dialog-title" style={{ color: "#ffffff", fontWeight: "600", padding: "15px", textUnderlineOffset: "3px", borderBottom: 'solid 1px #d2d6d5'}}>
      How Does This Work?
    </DialogTitle>
    <DialogContent style={{marginBottom: "0" }}>
      <DialogContentText 
        id="alert-dialog-description" 
        style={{ 
          color: '#ffffff', 
		  fontSize: '15px',
          marginBottom: '0',
          whiteSpace: 'pre-line', 
				
			padding: "19px"
        }}
      >
	  <div style={{textIndent: "30px"}}>
        This program takes the 19 source images, and turns them into an animated GIF by creating
        intermediate frames from scratch that transitions the images smoothly. The algorithm takes 
        each downloaded image from NASA, encodes it into GIF format, then decodes the data of that 
        GIF so that the images can be manipulated.
	  </div>
        <br />

	  <div style={{textIndent: "30px"}}>
        Every intermediate frame is produced by going through every pixel in each image, and finding 
        the absolute value of the difference between the first and second image’s RGB values (separately).
        The amount to increment per frame is then calculated based on how many intermediate frames 
        have been specified (I have chosen 80 for this site), and whether we need to go up or down 
        to reach the target value. Each pixel points to its own slice of interpolated pixels 
        (262144 total pixels in each 512x512 image, 80 total intermediate pixels per pixel index).
	  </div>	
        <br />
	
	  <div style={{textIndent: "30px"}}>
        Then, new images are created by stopping at the frame we are building, and constructing that 
        frame with every pixel that belongs to that frame’s index. Essentially, we are just stopping 
        at every X Y coordinate for a blank 512x512 image, 80 times, and adding the pixel that we’ve 
        calculated belongs in that place.
	  </div>
        <br />

	<div style={{textIndent: "30px"}}>
        This process repeats for all 19 images, until finally we turn all 1520 frames into one GIF.
        This is what you see!
	  </div>

	          <br />
		This type of algorithm is commonly referred to as <span><a target='_blank' href="https://en.wikipedia.org/wiki/Inbetweening" style={{textDecoration: "underline"}}>Inbetweening</a>.</span>
        <br />
	  <a target='_blank' href="https://github.com/beo1259/Sunscrape/blob/master/src/sunscrape.go" style={{textDecoration: "underline", color: "#91baff"}}>See the algorithm for yourself!</a>
      </DialogContentText>
    </DialogContent>
    <DialogActions style={{paddingTop: "13px",borderTop: 'solid 1px #d2d6d5'}}>
      <Button 
        onClick={handleClose}
        style={{ 
          textTransform: 'none', 
          fontSize: '16px',
	      border: "1px solid #cbd5e1",
		  color: "white",
		  borderRadius: "10px",
		  paddingRight:"6px",
		  paddingLeft:"6px",
				width: "100px",
		marginBottom: "6px",
		marginRight: "23px",
        }}
      >
        I Get It!
      </Button>
    </DialogActions>
  </Dialog>
</React.Fragment>

     );
}

