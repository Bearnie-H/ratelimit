# RateLimit
This package provides basic wrappers of io.Reader and io.Writer interfaces to allow limiting the read/write rate to no more than a specified number of bytes per second.

An situation where this could be helpful is in a server application, where the total network usage cannot exceed some threshold. The application could allocate and track blocks of read/write "rate" to go-routines throughout the application, and if all relevant read/write operations use the RateReader and RateWriter implementations, the applicaiton will never use more than the given pool of bandwidth.

## Read Example

``` go
// Read from the given reader, but don't try to go over 32KB/s read rate!
// Perhaps a network connection and we want to keep our footprint under this 32KB/s limit
// for complicated personal reasons.
func Foo(r io.Reader) error {

    // Wrap the given reader and limit it to only read at most 32KB/s
    R := NewRateReader(r, 1024*32)

    // Create a file on disk to write the contents to
    f, err := os.Create("file.txt")
    if err != nil {
        return err
    }
    defer f.Close()

    // Copy to the file on disk.
    if _, err := io.Copy(f, R); err != nil {
        return err
    }

    return nil
}
```

## Write Example
``` go
// Write out a file to a number of network connetions, but don't exceed a given write speed!
func Bar(f *os.File, writeLimit int, writers ...io.Writer) error {

    // Identify how many writers there are, to know how much speed to allocate
    // each one
    n := len(writers)
    for index := range writers {
        // Wrap each of the given Writers to be a RateWriter
        writers[index] = NewRateWriter(writers[index], writeLimit/n)
    }

    // Further encapsulate the set into a MultiWriter, to allow writing to all at once
    W := io.MultiWriter(writers...)

    // Finally, copy the file to all of the given writers
    if _, err := io.Copy(W, f); err != nil {
        return err
    }

    return nil
}
```