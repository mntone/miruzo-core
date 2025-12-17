#include <lcms2.h>

// Build (Ubuntu/Debian):
//   sudo apt install liblcms2-dev
//   gcc display_p3.c -o display_p3.o -llcms2
// Build (macOS/Homebrew):
//   brew install little-cms2
//   gcc display_p3.c -o display_p3.o -llcms2

int main(void) {
	// --- D65 white point ---
	cmsCIExyY whitePoint = {
		.x = 0.3127,
		.y = 0.3290,
		.Y = 1.0
	};

	// --- DCI-P3 primaries ---
	cmsCIExyYTRIPLE primaries = {
		.Red   = { 0.680, 0.320, 1.0 },
		.Green = { 0.265, 0.690, 1.0 },
		.Blue  = { 0.150, 0.060, 1.0 }
	};

	// --- sRGB TRC (EOTF: encoded -> linear) ---
	// parametric curve type 4
	// if x <= d: y = e*x + f
	// else:      y = (a*x + b)^g + c
	cmsFloat64Number srgbParams[7] = {
		2.4,                    // g
		1.0 / 1.055,            // a
		0.055 / 1.055,          // b
		0.0,                    // c
		0.04045,                // d
		1.0 / 12.92,            // e
		0.0                     // f
	};

	cmsToneCurve* trc =
		cmsBuildParametricToneCurve(NULL, 4, srgbParams);

	cmsToneCurve* trcs[3] = { trc, trc, trc };

	// --- Create RGB profile ---
	cmsHPROFILE profile =
		cmsCreateRGBProfile(&whitePoint, &primaries, trcs);

	// Set ICC profile version to v4 (important)
	cmsSetProfileVersion(profile, 4.3);

	// Description
	cmsWriteTag(profile, cmsSigProfileDescriptionTag,
		"Display P3 (DCI-P3 + D65 + sRGB TRC)");

	// --- Save ---
	cmsSaveProfileToFile(profile, "display_p3.icc");

	// --- Cleanup ---
	cmsFreeToneCurve(trc);
	cmsCloseProfile(profile);

	return 0;
}
