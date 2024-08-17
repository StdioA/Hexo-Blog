---
title: 一种 IPv6 地址编码方案
author: David Dai
tags:
  - IPv6
  - Unicode
  - UTF-8
categories:
  - 乱七八糟
date: 2022-06-10 20:51:14
toc: true
---

又搞了一些骚操作：把一个 IPv6 地址压缩成一个短字符串。

<!--more-->
# 背景
线上某张表有一个 `VARCHAR` 字段，用于存储 IP 地址。之前只存储 IPv4 地址，而 IPv4 地址的最大长度为 15（如 `255.255.255.255`），因此字段宽度只设置了 20。  
当我们要存储 IPv6 地址时，却发现 IPv6 地址的最大长度是 39（如 `1111:2222:3333:4444:5555:6666:7777:8888`），而变更字段类型的尝试也以失败告终，因此我们需要找到一种方法来将 v6 IP 塞进长度为 20 的 `VARCHAR` 字段中。

# 一些简单的尝试
随便找一个 v6 IP，如 `240e:17:ce8:fd00:52a8:6001:6e05:96f6`，然后尝试将它缩短。

* 去掉冒号是否可行？
 去掉冒号后还是有 32 位，不行。
* 将它变成二进制，然后再用 base64 编码？
 v6 IP 的（二进制）长度为 128 位，而 base64 一个字符可以存放 6 个二进制位，因此编码后字符串的长度至少有 128 / 6 = 21.3 位，再加上 base64 的固定 pad，最后需要 24 位。
 hmmmm…差一点点。
* 使用一些常见的压缩算法？
  IPv6 地址可以视为随机序列，本身压缩效率就不会很高，而且通常压缩后得到的都是二进制序列而不是合法字符串，感觉不太靠谱。
    
# 有没有压缩率更高的方案？
把 v6 IP 进行编码后，使它能够存入这个 `VARCHAR(20)` 的列，需要两个前提：
1. 编码后的字符串，应该是应该 UTF-8 编码下的合法字符串；
2. 字符串的长度需要小于等于 20.

MySQL 4 在计算 `VARCHAR` 的宽度时，是通过编码后的**字节数**来判断的，因此一个中文字符会占用 3 个宽度；而 MySQL 5 及之后的版本是通过**Unicode 字符（字元）数**来判断宽度的，这样，一个中文字符只占用 1 个宽度。

如果我们把 IP 地址变成一个二进制序列并分段，然后将每段序列变成一个 Unicode 字符，并通过 UTF-8 进行编码，是否可行？

查看 UTF-8 编码规则：

| 代码范围 | Unicode 标量值 | 二进制 UTF-8 格式 | 注释 | 能够存储的位（bit）数 |
| - | - | - | - | - |
| 000000 - 00007F | 0zzzzzzz | 0zzzzzzz（00-7F）|     ASCII 字元范围，字节由零开始 | 7 |
| 000080 - 0007FF | 00000yyy yyzzzzzz | 110yyyyy（C0-DF） 10zzzzzz（80-BF）| 第一个字节以110开始，之后的字节以10开始 | 11 |
| 000800 - 00D7FF <br/> 00E000 - 00FFFF | xxxxyyyy yyzzzzzz | 1110xxxx（E0-EF） 10yyyyyy 10zzzzzz |    第一个字节以1110开始，之后的字节以10开始 | 16 |
| 010000 - 10FFFF | 000wwwxx xxxxyyyy yyzzzzzz    | 11110www（F0-F7） 10xxxxxx 10yyyyyy 10zzzzzz | 由11110开始，之后的字节以10开始 | 21 |

MySQL 的 `utf8` <ruby>排序规则<rp>(</rp><rt>collation</rt><rp>)</rp></ruby>允许最多 3 个字节的 UTF-8 字符，而 `utf8mb4` 则能够支持 4 个字节的 UTF-8 字符。

虽然 UTF-8 4 字节可以容纳 21 位，但其中会触及未定义的 Unicode 平面，导致我们无法使用整个编码空间。因此我们可以尝试使用 16 位的长度对 IP 地址进行分段，然后将每段编码成占用 3 个字节（及以下）的字符。
一个 3 字节字符能够支持 16 位的 Unicode 标量，那么计算一下编码后的字符串长度：128 / 16 = 8，比现有的 `VARCHAR(20)` 字段宽度少很多，完全可行！


## 第一版实现：进制转换
编码时，将 IP 变为一个大整数，再转化为 `2^16` 进制数，每一位使用一个 Unicode 字符来表示；然后再通过 UTF-8 编码编入 `[]byte`，最后转换为 `string` 并返回。
解码时，读取每个 rune，获取它的 UTF-8 编码，再将其变为 Unicode 标量，再将 65536 进制数变成大整数，最后编码成 IP 地址。


举例：
> IP: 240e:17:ce8:fd00:52a8:6001:6e05:96f6
> 转为大整数：47924901830519682514395366120933810856
> 二进制： 10010000001110000000000001011…
> 进制转换切出第一段二进制序列 1001011011110110 -> 字符“零”
> 整段 IP 转码：␎೨ﴀ动态清零

实现如下：
```go
const padLength = 16

func IPv6ToUTF8(ip6 net.IP) string {
    var (
        builder    strings.Builder
        bytesArray [][]byte
    )
    // 1. 转换成 big.Int
    bigInt := big.NewInt(0).SetBytes(ip6.To16())
    // 2. 通过进制转换的方式，将大整数按 16 位切分
    mask := big.NewInt(1 << padLength) // 最多可以放 16 位
    for bigInt.Cmp(big.NewInt(0)) > 0 {
        // 获取当前段
        pad := big.NewInt(0).Mod(bigInt, mask).Uint64()
                
        // 3. 塞入 utf-8 中
        var numBytes []byte
        switch {
        case pad < 0x80: // 一字节字符
            numBytes = []byte{byte(pad)}
        case pad < 0x7FF: // 两字节字符
            var bytes = []byte{0b11000000, 0b10000000}
            bytes[1] += byte((pad & 0b111111))            // bit 0-5
            bytes[0] += byte((pad & (0b11111 << 6)) >> 6) // bit 6-10
            numBytes = bytes
        case pad <= 0xFFFF: // 三字节字符
            var bytes = []byte{0b11100000, 0b10000000, 0b10000000}
            bytes[2] += byte((pad & 0b111111))             // bit 0-5
            bytes[1] += byte((pad & (0b111111 << 6)) >> 6) // bit 6-11
            bytes[0] += byte((pad & (0b1111 << 12)) >> 12) // bit 12-15
            numBytes = bytes
        }
        bytesArray = append(bytesArray, numBytes)
        // 移位，处理下一段
        bigInt.Div(bigInt, mask)
    }
    // 之前是倒序放入 bytesArray 的，现在需要倒过来
    for i := len(bytesArray) - 1; i >= 0; i-- {
        builder.Write(bytesArray[i])
    }
    return builder.String()
}

func UTF8ToIPv6(s string) net.IP {
    var bitInt = big.NewInt(0)
    // 遍历 rune
    for _, r := range s {
        // 构建 unicode 字元，并完成 UTF-8 解码
        var num int64
        switch {
        case r < 0x80:
            num += int64(r)
        case r < 0x7ff:
            num += int64(r & 0b111111)              // bit 0-5
            num += int64(((r >> 6) & 0b11111) << 6) // bit 6-10
        case r <= 0xffff:
            num += int64(r & 0b111111)               // bit 0-5
            num += int64(((r >> 6) & 0b111111) << 6) // bit 6-11
            num += int64(((r >> 12) & 0b1111) << 12) // bit 12-15
        }
        // 将该段拼接到大整数中
        // bitInt = bitInt << padLength + num
        shifted := bitInt.Lsh(bitInt, padLength)
        bitInt.Add(shifted, big.NewInt(num))
    }
    // 把整数变回 IP
    ip := net.IP(bitInt.Bytes())
    return ip
}
```

## 第二版实现：bytes to rune
上面那段太复杂了，各种位运算乱七八糟，还要处理 UTF-8 的变长编码，有没有更简单的方案？
仔细看下可以发现，刚刚我们切段的时候，每段的长度是 16 位，刚好就是两个字节。因此，我们直接以 2 字节为单位，对 IP 进行切段，把每段变成 rune，再将 rune 拼成 string 即可。
恰巧，IPv6 地址分成了 8 段，每段也刚好是两字节。

举例：
> IP：240e:17:ce8:fd00:52a8:6001:6e05:96f6
> 转换为字节序列并按双字节切段: [\x24 \x0e] [\x00 \x17] [\x0c \xe8]....
> 将每段变为一个字符: \u240e \u0017 \u0e8  ....
> 最终字符串：␎೨ﴀ动态清零

重新实现后，代码更加简洁，同时编码效率提高了 17 倍，解码效率提高了 10 倍。

```go
func IPv6ToUTF8(ip6 net.IP) string {
    var (
        bytes   = []byte(ip6)       // 1. 转换成 bytes
        builder strings.Builder
    )
    builder.Grow(18)  // 存放 6 个字符，每个最多 3 字节，因此提前 grow
    for i := 0; i < len(bytes); i += 2 {     // 2. 按两字节分段
        r := rune(bytes[i])<<8 | rune(bytes[i+1])
        builder.WriteRune(r)   // 3. 将 rune 拼成 string
    }
    return builder.String()
}

func UTF8ToIPv6(ipStr string) net.IP {
    var (
        bytes = make([]byte, 16)
        i     int
    )
    // 把每个 rune（uint32）的低 16 位拆分成两个字节，并放到对应的位置
    for _, r := range ipStr {
        bytes[i] = byte(r >> 8)
        bytes[i+1] = byte(r & 0xff)
        i += 2
    }
    // 有了字节序列，转换为 IP 即可
    return net.IP(bytes)
}
```

## Unicode 代理对带来的小麻烦
仔细看下刚刚的 UTF-8 的编码表，可能会发现，三个字节可以容纳的字元范围是 `\u0800 - \ud7ff` 和 `\ue000 - \uffff`，不包含中间的  `\ud800 - \udfff` 区间。尝试了一下，如果字节段被映射到了这个区间内，它就会变成一个非法字符 "�"，而再解析出来的 Unicode 字元会变成 `\ufffd`，这样类似 `d800::dffe` 的 IP 在经过编码再解码后就会变成 `fffd::fffd`，这显然是不能接受的。

研究 Unicode 区段发现，这段字元处于“Unicode <ruby>代理对<rp>(</rp><rt>surrogate</rt><rp>)</rp></ruby>”区段，通常都会成对出现，单独出现时不会被视为一个合法的 Unicode 字符，因此在 UTF-8 解码时会被替换为一个“[替换字符](https://www.compart.com/en/unicode/U+FFFD)” `\ufffd`，也就是 "�"。所以，我们在编码时，如果遇到相关的字节范围，则需要避开这些区段。
一个简单的操作方式，是直接将其平移到 `\u1d800 - \u1dfff`。这段 Unicode 包含两个区段，分别为“[萨顿书写字母](https://www.compart.com/en/unicode/block/U+1D800)”和一个小的未定义区段，均为合法区段，因此可以解决非法区间问题。

不过，因为 `\u1d800` 在编码为 UTF-8 后会占用四个字节，因此需要保证数据库表使用的是 `utf8mb4` 排序规则. 如果数据库使用的是 `utf8` 排序规则，可以考虑缩短分段位数，比如只使用两字节的 UTF-8 编码，每段 11 位，这样需要 12 个字符就可以存下 IPv6 地址。

实现时，对相关区段进行特判即可：
```go
func IPv6ToUTF8(ip6 net.IP) string {
    var (
        bytes   = []byte(ip6)
        builder strings.Builder
    )
    builder.Grow(18)
    for i := 0; i < len(bytes); i += 2 {
        r := rune(bytes[i])<<8 | rune(bytes[i+1])
        // D800 - DFFF，涉及到三个代理对区段，无法编码，解码时会变成 FFFD
        // 所以需要手动更改区段到 1D800 - 1DFFF
        // 因此对应的 UTF-8 也会多一个字节，不过不影响 MySQL 字符长度
        if 0xD800 <= r && r <= 0xDFFF {
            r += 0x10000
        }
        builder.WriteRune(r)
    }
    return builder.String()
}

func UTF8ToIPv6(ipStr string) net.IP {
    var (
        bytes = make([]byte, 16)
        i     int
    )
    for _, r := range ipStr {
        // 这里碰到 1D800 - 1DFFF 可能会发生溢出，不过不影响计算，因为我们只需要中间的一字节，更高位的 1 可以丢弃
        bytes[i] = byte(r >> 8) 
        bytes[i+1] = byte(r & 0xff)
        i += 2
    }
    return net.IP(bytes)
}
```

# 收尾
完成了有效的 IPv6 压缩方案后，我们在存储 IP 时，直接原样存储 v4 IP；当存储 v6 IP 时，则使用上文提到的压缩方案，并在编码过后的字符串前添加一个 v4 IP 中不会出现的字符（如 ":"）作为前缀，用于区分 v4 和 v6 IP，即可得到一个完善的、同时兼容 v4 和 v6 IP 的存储方案。

```go
const ipv6Prefix = ":"

// EncodeIP 将 v6 IP 转换成压缩格式，v4 IP 原样返回
func EncodeIP(ipStr string) string {
    ip := net.ParseIP(ipStr)
    if ip.To4() == nil && ip.To16() != nil {
        // 压缩 IPv6
        return ipv6Prefix + IPv6ToUTF8(ip)
    }
    // 针对 IPv4 和非法 IP（如已经压缩过的 IP），直接原样返回
    return ipStr
}

func DecodeIP(encoded string) string {
    if strings.HasPrefix(encoded, ipv6Prefix) {
        return UTF8ToIPv6(encoded[1:]).String()
    }
    return encoded
}
```

# 参考资料
1. [Unicode Blocks](https://www.compart.com/en/unicode/block)
2. [UTF-8 - 维基百科](https://zh.m.wikipedia.org/zh-hans/UTF-8)
3. [Unicode 代理对](https://cloud.tencent.com/developer/article/1641938)
