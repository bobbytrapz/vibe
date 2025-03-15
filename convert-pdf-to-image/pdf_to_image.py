#!/usr/bin/env python3
"""
PDF to Image Converter

This script converts PDF files to image formats (PNG, JPEG, etc.)
"""

import argparse
import sys
from pathlib import Path
import fitz  # PyMuPDF
from PIL import Image


def convert_pdf_to_image(pdf_path: str | Path, output_dir: str | Path | None = None, format: str = "png", dpi: int = 300, pages: list[int] | None = None):
    """
    Convert a PDF file to image format
    
    Args:
        pdf_path (str): Path to the PDF file
        output_dir (str, optional): Directory to save images, defaults to same as PDF
        format (str, optional): Output image format (png, jpg, etc.), defaults to png
        dpi (int, optional): Image resolution, defaults to 300
        pages (list, optional): Specific pages to convert, defaults to all
    
    Returns:
        list: Paths to generated image files
    """
    # Convert string paths to Path objects
    if isinstance(pdf_path, str):
        pdf_path = Path(pdf_path)
    
    if not pdf_path.exists():
        print(f"Error: PDF file '{pdf_path}' not found")
        return []
    
    # Determine output directory
    if output_dir is None:
        output_dir = pdf_path.parent or Path(".")
    elif isinstance(output_dir, str):
        output_dir = Path(output_dir)
    
    # Create output directory if it doesn't exist
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # Get base filename without extension
    base_filename = pdf_path.stem
    
    # Open PDF
    try:
        pdf_document = fitz.open(str(pdf_path))
    except Exception as e:
        print(f"Error opening PDF: {e}")
        return []
    
    # Calculate zoom factor based on DPI (default PDF DPI is 72)
    zoom = dpi / 72
    
    generated_images = []
    
    # Process pages
    total_pages = len(pdf_document)
    
    # Determine which pages to convert
    if pages is None:
        pages_to_convert = range(total_pages)
    else:
        pages_to_convert = [p-1 for p in pages if 0 < p <= total_pages]  # Convert to 0-based index
    
    for page_num in pages_to_convert:
        print(f"Converting page {page_num + 1}/{total_pages}...")
        
        page = pdf_document.load_page(page_num)
        
        # Get page as a pixmap (image)
        mat = fitz.Matrix(zoom, zoom)
        pix = page.get_pixmap(matrix=mat, alpha=False)
        
        # Construct output filename
        if len(pages_to_convert) == 1 and total_pages == 1:
            # If there's only one page in the PDF and we're converting only that page
            output_filename = f"{base_filename}.{format}"
        else:
            # Multiple pages - add page number to filename
            output_filename = f"{base_filename}_page_{page_num + 1}.{format}"
        
        output_path = output_dir / output_filename
        
        # Save the pixmap
        pix.save(str(output_path))
        generated_images.append(str(output_path))
    
    pdf_document.close()
    print(f"Conversion complete. {len(generated_images)} images generated in {output_dir}")
    return generated_images


def main():
    parser = argparse.ArgumentParser(description="Convert PDF to image format")
    parser.add_argument("pdf_path", help="Path to the PDF file")
    parser.add_argument("-o", "--output-dir", help="Directory to save images (default: same as PDF)")
    parser.add_argument("-f", "--format", default="png", help="Output image format (png, jpg, etc.)")
    parser.add_argument("-d", "--dpi", type=int, default=300, help="Image resolution (DPI)")
    parser.add_argument("-p", "--pages", type=int, nargs="+", help="Specific pages to convert (e.g., -p 1 3 5)")
    
    args = parser.parse_args()
    
    try:
        convert_pdf_to_image(
            args.pdf_path, 
            args.output_dir, 
            args.format, 
            args.dpi, 
            args.pages
        )
    except Exception as e:
        print(f"Error during conversion: {e}")
        return 1
    
    return 0


if __name__ == "__main__":
    sys.exit(main()) 
